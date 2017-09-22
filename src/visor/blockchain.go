package visor

import (
	"errors"
	"fmt"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/visor/blockdb"

	"bytes"
)

var (
	// DebugLevel1 checks for extremely unlikely conditions (10e-40)
	DebugLevel1 = true
	// DebugLevel2 enable checks for impossible conditions
	DebugLevel2 = true

	// ErrUnspentNotExist represents the error of unspent output in a tx does not exist
	ErrUnspentNotExist = errors.New("Unspent output does not exist")
	// ErrSignatureLost signature lost error
	ErrSignatureLost = errors.New("signature lost")
)

const (
	// SigVerifyTheadNum  signature verifycation goroutine number
	SigVerifyTheadNum = 4
)

//Warning: 10e6 is 10 million, 1e6 is 1 million

// Note: DebugLevel1 adds additional checks for hash collisions that
// are unlikely to occur. DebugLevel2 adds checks for conditions that
// can only occur through programmer error and malice.

// Note: a droplet is the base coin unit. Each Skycoin is one million droplets

//Termonology:
// UXTO - unspent transaction outputs
// UX - outputs10
// TX - transactions

//Notes:
// transactions (TX) consume outputs (UX) and produce new outputs (UX)
// Tx.Uxi() - set of outputs consumed by transaction
// Tx.Uxo() - set of outputs created by transaction

// BlockTree provide method for access
type BlockTree interface {
	AddBlock(b *coin.Block) error
	AddBlockWithTx(tx *bolt.Tx, b *coin.Block) error
	RemoveBlock(b *coin.Block) error
	GetBlock(hash cipher.SHA256) *coin.Block
	GetBlockInDepth(dep uint64, filter func(hps []coin.HashPair) cipher.SHA256) *coin.Block
}

// Walker function for go through blockchain
type Walker func(hps []coin.HashPair) cipher.SHA256

// BlockListener notify the register when new block is appended to the chain
type BlockListener func(b coin.Block)

// Blockchain maintains blockchain and provides apis for accessing the chain.
type Blockchain struct {
	db          *bolt.DB
	tree        BlockTree
	pubkey      cipher.PubKey
	walker      Walker
	blkListener []BlockListener

	// arbitrating mode, if in arbitrating mode, when master node execute blocks,
	// the invalid transaction will be skipped and continue the next; otherwise,
	// node will throw the error and return.
	arbitrating bool
	chain       *blockdb.Blockchain
	sigs        *blockdb.BlockSigs
}

// Option represents the option when creating the blockchain
type Option func(*Blockchain)

// DefaultWalker default blockchain walker
func DefaultWalker(hps []coin.HashPair) cipher.SHA256 {
	return hps[0].Hash
}

// NewBlockchain use the walker go through the tree and update the head and unspent outputs.
func NewBlockchain(db *bolt.DB, pubkey cipher.PubKey, ops ...Option) (*Blockchain, error) {
	// creates blockchain tree
	tree, err := blockdb.NewBlockTree(db)
	if err != nil {
		return nil, err
	}

	chainstore, err := blockdb.NewBlockchain(db)
	if err != nil {
		return nil, err
	}

	sigs, err := blockdb.NewBlockSigs(db)
	if err != nil {
		return nil, err
	}

	bc := &Blockchain{
		db:     db,
		pubkey: pubkey,
		tree:   tree,
		walker: DefaultWalker,
		chain:  chainstore,
		sigs:   sigs,
	}

	for _, op := range ops {
		op(bc)
	}

	if err := bc.walkTree(); err != nil {
		return nil, err
	}

	// verify signature
	if err := bc.VerifySigs(); err != nil {
		return nil, err
	}

	return bc, nil
}

// Arbitrating option to change the mode
func Arbitrating(enable bool) Option {
	return func(bc *Blockchain) {
		bc.arbitrating = enable
	}
}

func BlockchainWalker(walker Walker) Option {
	return func(bc *Blockchain) {
		bc.walker = walker
	}
}

func (bc *Blockchain) walkTree() error {
	var dep uint64
	var preBlock *coin.Block
	head := bc.Head()
	if head != nil {
		dep = head.Seq() + 1
	}

	preBlock = head

	for {
		b := bc.tree.GetBlockInDepth(dep, bc.walker)
		if b == nil {
			break
		}
		if dep > 0 {
			if b.PreHashHeader() != preBlock.HashHeader() {
				return errors.New("walk tree failed, pre hash header not match")
			}
		}
		preBlock = b
		if err := bc.updateUnspent(*b); err != nil {
			return fmt.Errorf("update unspent failed, err: %v", err.Error())
		}

		dep++
	}
	return nil
}

func (bc *Blockchain) headSeq() int64 {
	return bc.chain.HeadSeq()
}

func (bc *Blockchain) processBlockWithTx(tx *bolt.Tx, b coin.Block) (coin.Block, error) {
	if bc.Len() > 0 && !bc.isGenesisBlock(b) {
		if err := bc.verifyBlockHeader(b); err != nil {
			return coin.Block{}, err
		}
		txns, err := bc.processTransactions(b.Body.Transactions)
		if err != nil {
			return coin.Block{}, err
		}
		b.Body.Transactions = txns

		if err := bc.verifyUxHash(b); err != nil {
			return coin.Block{}, err
		}
	}

	if err := bc.chain.ProcessBlockWithTx(tx, &b); err != nil {
		return coin.Block{}, err
	}

	return b, nil
}

func (bc *Blockchain) processBlock(b coin.Block) (coin.Block, error) {
	if bc.Len() > 0 && !bc.isGenesisBlock(b) {
		if err := bc.verifyBlockHeader(b); err != nil {
			return coin.Block{}, err
		}
		txns, err := bc.processTransactions(b.Body.Transactions)
		if err != nil {
			return coin.Block{}, err
		}
		b.Body.Transactions = txns

		if err := bc.verifyUxHash(b); err != nil {
			return coin.Block{}, err
		}
	}

	if err := bc.chain.ProcessBlock(&b); err != nil {
		return coin.Block{}, err
	}

	return b, nil
}

// Unspent returns the unspent outputs pool
func (bc *Blockchain) Unspent() *blockdb.UnspentPool {
	return bc.chain.Unspent
}

// Len returns the length of current blockchain.
func (bc Blockchain) Len() uint64 {
	head := bc.Head()
	if head != nil {
		return head.Seq() + 1
	}
	return 0
}

// GetGenesisBlock get genesis block.
func (bc Blockchain) GetGenesisBlock() *coin.Block {
	return bc.tree.GetBlockInDepth(0, bc.walker)
}

// CreateGenesisBlock creates genesis block in blockchain.
func (bc *Blockchain) CreateGenesisBlock(genesisAddr cipher.Address, genesisCoins, timestamp uint64) (coin.Block, error) {
	txn := coin.Transaction{}
	txn.PushOutput(genesisAddr, genesisCoins, genesisCoins)
	body := coin.BlockBody{Transactions: coin.Transactions{txn}}
	prevHash := cipher.SHA256{}
	head := coin.BlockHeader{
		Time:     timestamp,
		BodyHash: body.Hash(),
		PrevHash: prevHash,
		BkSeq:    0,
		Version:  0,
		Fee:      0,
		UxHash:   cipher.SHA256{},
	}
	b := coin.Block{
		Head: head,
		Body: body,
	}

	if err := bc.db.Update(func(tx *bolt.Tx) error {
		var er error
		b, er = bc.processBlockWithTx(tx, b)
		if er != nil {
			return er
		}

		return bc.addBlockWithTx(tx, &b)
	}); err != nil {
		return coin.Block{}, err
	}

	bc.notify(b)
	return b, nil
}

// AddSig adds block signature
func (bc *Blockchain) AddSig(sb *coin.SignedBlock) error {
	return bc.sigs.Add(sb)
}

func (bc *Blockchain) addBlock(b *coin.Block) error {
	return bc.tree.AddBlock(b)
}

func (bc *Blockchain) addBlockWithTx(tx *bolt.Tx, b *coin.Block) error {
	return bc.tree.AddBlockWithTx(tx, b)
}

// GetBlock get block of specific hash in the blockchain, return nil on not found.
func (bc Blockchain) GetBlock(hash cipher.SHA256) *coin.Block {
	return bc.tree.GetBlock(hash)
}

// Head returns the most recent confirmed block
func (bc Blockchain) Head() *coin.Block {
	headSeq := bc.headSeq()
	if headSeq < 0 {
		return nil
	}

	return bc.GetBlockInDepth(uint64(headSeq))
}

// Time returns time of last block
// used as system clock indepedent clock for coin hour calculations
// TODO: Deprecate
func (bc *Blockchain) Time() uint64 {
	return bc.Head().Time()
}

// NewBlockFromTransactions creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification.  Transactions must be sorted.
func (bc Blockchain) NewBlockFromTransactions(txns coin.Transactions, currentTime uint64) (*coin.Block, error) {
	if currentTime <= bc.Time() {
		return nil, errors.New("Time can only move forward")
	}

	if len(txns) == 0 {
		return nil, errors.New("No transactions")
	}
	txns, err := bc.processTransactions(txns)
	if err != nil {
		return nil, err
	}
	uxHash := bc.Unspent().GetUxHash()

	b, err := coin.NewBlock(*bc.Head(), currentTime, uxHash, txns, bc.TransactionFee)
	if err != nil {
		return nil, err
	}

	//make sure block is valid
	if DebugLevel2 == true {
		if err := bc.verifyBlockHeader(*b); err != nil {
			return nil, err
		}
		txns, err := bc.processTransactions(b.Body.Transactions)
		if err != nil {
			logger.Panic("Impossible Error: not allowed to fail")
		}
		b.Body.Transactions = txns
	}
	return b, nil
}

// ExecuteBlockWithTx attempts to append block to blockchain with *bolt.Tx
func (bc *Blockchain) ExecuteBlockWithTx(tx *bolt.Tx, sb *coin.SignedBlock) error {
	if err := bc.sigs.AddWithTx(tx, sb); err != nil {
		return err
	}

	b := sb.Block
	b.Head.PrevHash = bc.Head().HashHeader()
	nb, err := bc.processBlockWithTx(tx, b)
	if err != nil {
		return err
	}

	if err := bc.addBlockWithTx(tx, &nb); err != nil {
		return err
	}

	bc.notify(nb)
	return nil
}

func (bc *Blockchain) updateUnspent(b coin.Block) error {
	_, err := bc.processBlock(b)
	return err
}

// isGenesisBlock checks if the block is genesis block
func (bc Blockchain) isGenesisBlock(b coin.Block) bool {
	return bc.GetGenesisBlock().HashHeader() == b.HashHeader()
}

// Compares the state of the current UxHash hash to state of unspent
// output pool.
func (bc Blockchain) verifyUxHash(b coin.Block) error {
	uxHash := bc.Unspent().GetUxHash()

	if !bytes.Equal(b.Head.UxHash[:], uxHash[:]) {
		return errors.New("UxHash does not match")
	}
	return nil
}

// VerifyTransaction checks that the inputs to the transaction exist,
// that the transaction does not create or destroy coins and that the
// signatures on the transaction are valid
func (bc Blockchain) VerifyTransaction(tx coin.Transaction) error {
	//CHECKLIST: DONE: check for duplicate ux inputs/double spending
	//CHECKLIST: DONE: check that inputs of transaction have not been spent
	//CHECKLIST: DONE: check there are no duplicate outputs

	// Q: why are coin hours based on last block time and not
	// current time?
	// A: no two computers will agree on system time. Need system clock
	// indepedent timing that everyone agrees on. fee values would depend on
	// local clock

	// Check transaction type and length
	// Check for duplicate outputs
	// Check for duplicate inputs
	// Check for invalid hash
	// Check for no inputs
	// Check for no outputs
	// Check for non 1e6 multiple coin outputs
	// Check for zero coin outputs
	// Check valid looking signatures
	if err := tx.Verify(); err != nil {
		return err
	}

	uxIn, err := bc.Unspent().GetArray(tx.In)
	if err != nil {
		return err
	}
	// Checks whether ux inputs exist,
	// Check that signatures are allowed to spend inputs
	if err := tx.VerifyInput(uxIn); err != nil {
		return err
	}

	// Get the UxOuts we expect to have when the block is created.
	uxOut := coin.CreateUnspents(bc.Head().Head, tx)
	// Check that there are any duplicates within this set
	if uxOut.HasDupes() {
		return errors.New("Duplicate unspent outputs in transaction")
	}
	if DebugLevel1 {
		// Check that new unspents don't collide with existing.  This should
		// also be checked in verifyTransactions
		for i := range uxOut {
			if bc.Unspent().Contains(uxOut[i].Hash()) {
				return errors.New("New unspent collides with existing unspent")
			}
		}
	}

	// Check that no coins are lost, and sufficient coins and hours are spent
	err = coin.VerifyTransactionSpending(bc.Time(), uxIn, uxOut)
	if err != nil {
		return err
	}
	return nil
}

// GetBlockInDepth return block whose BkSeq is seq.
func (bc Blockchain) GetBlockInDepth(dep uint64) *coin.Block {
	return bc.tree.GetBlockInDepth(dep, bc.walker)
}

// GetBlocks return blocks whose seq are in the range of start and end.
func (bc Blockchain) GetBlocks(start, end uint64) []coin.Block {
	if start > end {
		return []coin.Block{}
	}

	blocks := []coin.Block{}
	for i := start; i <= end; i++ {
		b := bc.tree.GetBlockInDepth(i, bc.walker)
		if b == nil {
			break
		}
		blocks = append(blocks, *b)
	}
	return blocks
}

// GetSig returns signautre of given block
func (bc Blockchain) GetSig(hash cipher.SHA256) (cipher.Sig, bool, error) {
	return bc.sigs.Get(hash)
}

// GetLastBlocks return the latest N blocks.
func (bc Blockchain) GetLastBlocks(num uint64) []coin.Block {
	var blocks []coin.Block
	if num == 0 {
		return blocks
	}

	end := bc.Head().Seq()
	start := end - num + 1
	if start < 0 {
		start = 0
	}
	return bc.GetBlocks(start, end)
}

/* Private */

// Validates a set of Transactions, individually, against each other and
// against the Blockchain.  If firstFail is true, it will return an error
// as soon as it encounters one.  Else, it will return an array of
// Transactions that are valid as a whole.  It may return an error if
// firstFalse is false, if there is no way to filter the txns into a valid
// array, i.e. processTransactions(processTransactions(txn, false), true)
// should not result in an error, unless all txns are invalid.
// TODO:
//  - move arbitration to visor
//  - blockchain should have strict checking
func (bc Blockchain) processTransactions(txs coin.Transactions) (coin.Transactions, error) {
	// copy txs so that the following code won't modify the origianl txs
	txns := make(coin.Transactions, len(txs))
	copy(txns, txs)

	// Transactions need to be sorted by fee and hash before arbitrating
	if bc.arbitrating {
		txns = coin.SortTransactions(txns, bc.TransactionFee)
	}
	//TODO: audit
	if len(txns) == 0 {
		if bc.arbitrating {
			return txns, nil
		}
		// If there are no transactions, a block should not be made
		return nil, errors.New("No transactions")
	}

	skip := make(map[int]struct{})
	uxHashes := make(coin.UxHashSet, len(txns))
	for i, tx := range txns {
		// Check the transaction against itself.  This covers the hash,
		// signature indices and duplicate spends within itself
		err := bc.VerifyTransaction(tx)
		if err != nil {
			if bc.arbitrating {
				skip[i] = struct{}{}
				continue
			} else {
				return nil, err
			}
		}

		// Check that each pending unspent will be unique
		uxb := coin.UxBody{
			SrcTransaction: tx.Hash(),
		}
		for _, to := range tx.Out {
			uxb.Coins = to.Coins
			uxb.Hours = to.Hours
			uxb.Address = to.Address
			h := uxb.Hash()
			_, exists := uxHashes[h]
			if exists {
				if bc.arbitrating {
					skip[i] = struct{}{}
					continue
				} else {
					m := "Duplicate unspent output across transactions"
					return nil, errors.New(m)
				}
			}
			if DebugLevel1 {
				// Check that the expected unspent is not already in the pool.
				// This should never happen because its a hash collision
				if bc.Unspent().Contains(h) {
					if bc.arbitrating {
						skip[i] = struct{}{}
						continue
					} else {
						m := "Output hash is in the UnspentPool"
						return nil, errors.New(m)
					}
				}
			}
			uxHashes[h] = byte(1)
		}
	}

	// Filter invalid transactions before arbitrating between colliding ones
	if len(skip) > 0 {
		newtxns := make(coin.Transactions, len(txns)-len(skip))
		j := 0
		for i := range txns {
			if _, shouldSkip := skip[i]; !shouldSkip {
				newtxns[j] = txns[i]
				j++
			}
		}
		txns = newtxns
		skip = make(map[int]struct{})
	}

	// Check to ensure that there are no duplicate spends in the entire block,
	// and that we aren't creating duplicate outputs.  Duplicate outputs
	// within a single Transaction are already checked by VerifyTransaction
	hashes := txns.Hashes()
	for i := 0; i < len(txns)-1; i++ {
		s := txns[i]
		for j := i + 1; j < len(txns); j++ {
			t := txns[j]
			if DebugLevel1 {
				if hashes[i] == hashes[j] {
					// This is a non-recoverable error for filtering, and
					// should never occur.  It indicates a hash collision
					// amongst different txns. Duplicate transactions are
					// caught earlier, when duplicate expected outputs are
					// checked for, and will not trigger this.
					return nil, errors.New("Duplicate transaction")
				}
			}
			for a := range s.In {
				for b := range t.In {
					if s.In[a] == t.In[b] {
						if bc.arbitrating {
							// The txn with the highest fee and lowest hash
							// is chosen when attempting a double spend.
							// Since the txns are sorted, we skip the 2nd
							// iterable
							skip[j] = struct{}{}
						} else {
							m := "Cannot spend output twice in the same block"
							return nil, errors.New(m)
						}
					}
				}
			}
		}
	}

	// Filter the final results, if necessary
	if len(skip) > 0 {
		newtxns := make(coin.Transactions, 0, len(txns)-len(skip))
		for i := range txns {
			if _, shouldSkip := skip[i]; !shouldSkip {
				newtxns = append(newtxns, txns[i])
			}
		}
		return newtxns, nil
	}

	return txns, nil
}

// TransactionFee calculates the current transaction fee in coinhours of a Transaction
func (bc Blockchain) TransactionFee(t *coin.Transaction) (uint64, error) {
	headTime := bc.Time()
	inUxs, err := bc.Unspent().GetArray(t.In)
	if err != nil {
		return 0, err
	}

	return TransactionFee(t, headTime, inUxs)
}

// VerifySigs checks that BlockSigs state correspond with coin.Blockchain state
// and that all signatures are valid.
func (bc *Blockchain) VerifySigs() error {
	head := bc.Head()
	if head == nil {
		return nil
	}

	seqC := make(chan uint64)

	shutdown, errC := bc.sigVerifier(seqC)

	for i := uint64(0); i <= head.Seq(); i++ {
		seqC <- i
	}

	shutdown()

	return <-errC
}

// signature verifier will get block seq from seqC channel,
// and have multiple thread to do signature verification.
func (bc *Blockchain) sigVerifier(seqC chan uint64) (func(), <-chan error) {
	quitC := make(chan struct{})
	wg := sync.WaitGroup{}
	errC := make(chan error, 1)
	for i := 0; i < SigVerifyTheadNum; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case seq := <-seqC:
					if err := bc.verifyBlockSig(seq); err != nil {
						errC <- err
						return
					}
				case <-quitC:
					return
				}
			}
		}(i)
	}

	return func() {
		close(quitC)
		wg.Wait()
		select {
		case errC <- nil:
			// no error
		default:
			// already has error in errC
		}
	}, errC
}

func (bc *Blockchain) getSignedBlock(seq uint64) (*coin.SignedBlock, bool, error) {
	b := bc.GetBlockInDepth(seq)
	if b == nil {
		return nil, false, nil
	}

	// get sig
	sig, exist, err := bc.sigs.Get(b.HashHeader())
	if err != nil {
		return nil, false, fmt.Errorf("get signature of block %d failed: %v", seq, err)
	}

	if !exist {
		return nil, false, ErrSignatureLost
	}

	sb := &coin.SignedBlock{
		Block: *b,
		Sig:   sig,
	}

	return sb, true, nil
}

func (bc *Blockchain) verifyBlockSig(seq uint64) error {
	sb, exist, err := bc.getSignedBlock(seq)
	if err != nil {
		return err
	}

	if !exist {
		return fmt.Errorf("found no block in seq %d", seq)
	}

	return cipher.VerifySignature(bc.pubkey, sb.Sig, sb.Block.HashHeader())
}

// VerifyBlockHeader Returns error if the BlockHeader is not valid
func (bc Blockchain) verifyBlockHeader(b coin.Block) error {
	//check BkSeq
	head := bc.Head()
	if b.Head.BkSeq != head.Head.BkSeq+1 {
		return errors.New("BkSeq invalid")
	}
	//check Time, only requirement is that its monotonely increasing
	if b.Head.Time <= head.Head.Time {
		return errors.New("Block time must be > head time")
	}
	// Check block hash against previous head
	if b.Head.PrevHash != head.HashHeader() {
		return errors.New("PrevHash does not match current head")
	}
	if b.HashBody() != b.Head.BodyHash {
		return errors.New("Computed body hash does not match")
	}
	return nil
}

// BindListener register the listener to blockchain, when new block appended, the listener will be invoked.
func (bc *Blockchain) BindListener(ls BlockListener) {
	bc.blkListener = append(bc.blkListener, ls)
}

// notifies the listener the new block.
func (bc *Blockchain) notify(b coin.Block) {
	for _, l := range bc.blkListener {
		l(b)
	}
}
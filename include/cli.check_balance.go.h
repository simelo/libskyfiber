typedef struct {
    GoString_ Coins;
    GoString_ Hours;
} cli__Balance;
typedef struct {
    cli__Balance Confirmed;
    cli__Balance Spendable;
    cli__Balance Expected;
    GoString_ Address;
} cli__AddressBalances;
typedef struct {
    cli__Balance Confirmed;
    cli__Balance Spendable;
    cli__Balance Expected;
    GoSlice_ Addresses;
} cli__BalanceResult;
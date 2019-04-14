#include "test_main.h"
#include <check.h>
// run suite
int main(void) {
  int number_failed = 0;
  int number_failed_fork = 0;
  SRunner *sr = srunner_create(cipher_address());
  SRunner *sr_fork = srunner_create(coin_transaction_fork());
  srunner_add_suite(sr, cipher_bitcoin());                         // ok
  srunner_add_suite(sr, cipher_crypto());                          // ok
  srunner_add_suite(sr, cipher_encrypt_scrypt_chacha20poly1305()); // ok
  srunner_add_suite(sr, cipher_hash());                            // ok
  srunner_add_suite(sr, coin_blocks());                            // ok
  srunner_add_suite(sr, coin_coin());                              // ok
  srunner_add_suite(sr, coin_math());                              // ok
  srunner_add_suite(sr, coin_output());                            // ok
  srunner_add_suite(sr, coin_transaction());                       // ok
  srunner_add_suite(sr, param_distribution());                     // ok
  srunner_add_suite(sr, util_droplet());                           // ok
  srunner_add_suite(sr, util_fee());                               // ok
  srunner_add_suite(sr, cipher_testsuite());
  srunner_set_fork_status(sr, CK_NOFORK);
  srunner_set_fork_status(sr_fork, CK_FORK);
  srunner_run_all(sr, CK_VERBOSE);
  srunner_run_all(sr_fork, CK_VERBOSE);
  number_failed = srunner_ntests_failed(sr);
  number_failed_fork = srunner_ntests_failed(sr_fork);
  srunner_free(sr);
  srunner_free(sr_fork);
  sr = NULL;
  sr_fork = NULL;
  return (number_failed == 0) ? EXIT_SUCCESS : EXIT_FAILURE;
}
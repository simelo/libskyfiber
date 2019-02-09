# coding: utf-8

"""
    Skycoin REST API.

    Skycoin is a next-generation cryptocurrency.  # noqa: E501

    OpenAPI spec version: 0.25.0
    Contact: skycoin.doe@example.com
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


from __future__ import absolute_import

import unittest

import swagger_client
from swagger_client.api.default_api import DefaultApi  # noqa: E501
from swagger_client.rest import ApiException


class TestDefaultApi(unittest.TestCase):
    """DefaultApi unit test stubs"""

    def setUp(self):
        self.api = swagger_client.api.default_api.DefaultApi()  # noqa: E501

    def tearDown(self):
        pass

    def test_address_count(self):
        """Test case for address_count

        Returns the total number of unique address that have coins.  # noqa: E501
        """
        pass

    def test_coin_supply(self):
        """Test case for coin_supply

        """
        pass

    def test_csrf(self):
        """Test case for csrf

        Creates a new CSRF token. Previous CSRF tokens are invalidated by this call.  # noqa: E501
        """
        pass

    def test_default_connections(self):
        """Test case for default_connections

        defaultConnectionsHandler returns the list of default hardcoded bootstrap addresses.\\n They are not necessarily connected to.  # noqa: E501
        """
        pass

    def test_health(self):
        """Test case for health

        Returns node health data.  # noqa: E501
        """
        pass

    def test_network_connection(self):
        """Test case for network_connection

        This endpoint returns a specific connection.  # noqa: E501
        """
        pass

    def test_network_connections(self):
        """Test case for network_connections

        This endpoint returns all outgoings connections.  # noqa: E501
        """
        pass

    def test_network_connections_disconnect(self):
        """Test case for network_connections_disconnect

        """
        pass

    def test_network_connections_exchange(self):
        """Test case for network_connections_exchange

        """
        pass

    def test_network_connections_trust(self):
        """Test case for network_connections_trust

        trustConnectionsHandler returns all trusted connections.\\n They are not necessarily connected to. In the default configuration, these will be a subset of the default hardcoded bootstrap addresses.  # noqa: E501
        """
        pass

    def test_resend_unconfirmed_txns(self):
        """Test case for resend_unconfirmed_txns

        """
        pass

    def test_verify_address(self):
        """Test case for verify_address

        healthHandler returns node health data.  # noqa: E501
        """
        pass

    def test_version(self):
        """Test case for version

        """
        pass

    def test_wallet(self):
        """Test case for wallet

        Returns a wallet by id.  # noqa: E501
        """
        pass

    def test_wallet_balance(self):
        """Test case for wallet_balance

        Returns the wallet's balance, both confirmed and predicted.  The predicted balance is the confirmed balance minus the pending spends.  # noqa: E501
        """
        pass

    def test_wallet_folder(self):
        """Test case for wallet_folder

        """
        pass

    def test_wallet_new_address(self):
        """Test case for wallet_new_address

        """
        pass

    def test_wallet_new_seed(self):
        """Test case for wallet_new_seed

        """
        pass

    def test_wallet_seed(self):
        """Test case for wallet_seed

        This endpoint only works for encrypted wallets. If the wallet is unencrypted, The seed will be not returned.  # noqa: E501
        """
        pass

    def test_wallet_update(self):
        """Test case for wallet_update

        Update the wallet.  # noqa: E501
        """
        pass

    def test_wallets(self):
        """Test case for wallets

        """
        pass


if __name__ == '__main__':
    unittest.main()
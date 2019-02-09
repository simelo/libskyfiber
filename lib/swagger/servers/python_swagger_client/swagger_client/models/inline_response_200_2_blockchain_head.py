# coding: utf-8

"""
    Skycoin REST API.

    Skycoin is a next-generation cryptocurrency.  # noqa: E501

    OpenAPI spec version: 0.25.0
    Contact: skycoin.doe@example.com
    Generated by: https://github.com/swagger-api/swagger-codegen.git
"""


import pprint
import re  # noqa: F401

import six


class InlineResponse2002BlockchainHead(object):
    """NOTE: This class is auto generated by the swagger code generator program.

    Do not edit the class manually.
    """

    """
    Attributes:
      swagger_types (dict): The key is attribute name
                            and the value is attribute type.
      attribute_map (dict): The key is attribute name
                            and the value is json key in definition.
    """
    swagger_types = {
        'block_hash': 'str',
        'fee': 'int',
        'previous_block_hash': 'str',
        'seq': 'str',
        'timestamp': 'int',
        'tx_body_hash': 'str',
        'ux_hash': 'str',
        'version': 'int'
    }

    attribute_map = {
        'block_hash': 'block_hash',
        'fee': 'fee',
        'previous_block_hash': 'previous_block_hash',
        'seq': 'seq',
        'timestamp': 'timestamp',
        'tx_body_hash': 'tx_body_hash',
        'ux_hash': 'ux_hash',
        'version': 'version'
    }

    def __init__(self, block_hash=None, fee=None, previous_block_hash=None, seq=None, timestamp=None, tx_body_hash=None, ux_hash=None, version=None):  # noqa: E501
        """InlineResponse2002BlockchainHead - a model defined in Swagger"""  # noqa: E501

        self._block_hash = None
        self._fee = None
        self._previous_block_hash = None
        self._seq = None
        self._timestamp = None
        self._tx_body_hash = None
        self._ux_hash = None
        self._version = None
        self.discriminator = None

        if block_hash is not None:
            self.block_hash = block_hash
        if fee is not None:
            self.fee = fee
        if previous_block_hash is not None:
            self.previous_block_hash = previous_block_hash
        if seq is not None:
            self.seq = seq
        if timestamp is not None:
            self.timestamp = timestamp
        if tx_body_hash is not None:
            self.tx_body_hash = tx_body_hash
        if ux_hash is not None:
            self.ux_hash = ux_hash
        if version is not None:
            self.version = version

    @property
    def block_hash(self):
        """Gets the block_hash of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The block_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: str
        """
        return self._block_hash

    @block_hash.setter
    def block_hash(self, block_hash):
        """Sets the block_hash of this InlineResponse2002BlockchainHead.


        :param block_hash: The block_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: str
        """

        self._block_hash = block_hash

    @property
    def fee(self):
        """Gets the fee of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The fee of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: int
        """
        return self._fee

    @fee.setter
    def fee(self, fee):
        """Sets the fee of this InlineResponse2002BlockchainHead.


        :param fee: The fee of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: int
        """

        self._fee = fee

    @property
    def previous_block_hash(self):
        """Gets the previous_block_hash of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The previous_block_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: str
        """
        return self._previous_block_hash

    @previous_block_hash.setter
    def previous_block_hash(self, previous_block_hash):
        """Sets the previous_block_hash of this InlineResponse2002BlockchainHead.


        :param previous_block_hash: The previous_block_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: str
        """

        self._previous_block_hash = previous_block_hash

    @property
    def seq(self):
        """Gets the seq of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The seq of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: str
        """
        return self._seq

    @seq.setter
    def seq(self, seq):
        """Sets the seq of this InlineResponse2002BlockchainHead.


        :param seq: The seq of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: str
        """

        self._seq = seq

    @property
    def timestamp(self):
        """Gets the timestamp of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The timestamp of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: int
        """
        return self._timestamp

    @timestamp.setter
    def timestamp(self, timestamp):
        """Sets the timestamp of this InlineResponse2002BlockchainHead.


        :param timestamp: The timestamp of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: int
        """

        self._timestamp = timestamp

    @property
    def tx_body_hash(self):
        """Gets the tx_body_hash of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The tx_body_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: str
        """
        return self._tx_body_hash

    @tx_body_hash.setter
    def tx_body_hash(self, tx_body_hash):
        """Sets the tx_body_hash of this InlineResponse2002BlockchainHead.


        :param tx_body_hash: The tx_body_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: str
        """

        self._tx_body_hash = tx_body_hash

    @property
    def ux_hash(self):
        """Gets the ux_hash of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The ux_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: str
        """
        return self._ux_hash

    @ux_hash.setter
    def ux_hash(self, ux_hash):
        """Sets the ux_hash of this InlineResponse2002BlockchainHead.


        :param ux_hash: The ux_hash of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: str
        """

        self._ux_hash = ux_hash

    @property
    def version(self):
        """Gets the version of this InlineResponse2002BlockchainHead.  # noqa: E501


        :return: The version of this InlineResponse2002BlockchainHead.  # noqa: E501
        :rtype: int
        """
        return self._version

    @version.setter
    def version(self, version):
        """Sets the version of this InlineResponse2002BlockchainHead.


        :param version: The version of this InlineResponse2002BlockchainHead.  # noqa: E501
        :type: int
        """

        self._version = version

    def to_dict(self):
        """Returns the model properties as a dict"""
        result = {}

        for attr, _ in six.iteritems(self.swagger_types):
            value = getattr(self, attr)
            if isinstance(value, list):
                result[attr] = list(map(
                    lambda x: x.to_dict() if hasattr(x, "to_dict") else x,
                    value
                ))
            elif hasattr(value, "to_dict"):
                result[attr] = value.to_dict()
            elif isinstance(value, dict):
                result[attr] = dict(map(
                    lambda item: (item[0], item[1].to_dict())
                    if hasattr(item[1], "to_dict") else item,
                    value.items()
                ))
            else:
                result[attr] = value

        return result

    def to_str(self):
        """Returns the string representation of the model"""
        return pprint.pformat(self.to_dict())

    def __repr__(self):
        """For `print` and `pprint`"""
        return self.to_str()

    def __eq__(self, other):
        """Returns true if both objects are equal"""
        if not isinstance(other, InlineResponse2002BlockchainHead):
            return False

        return self.__dict__ == other.__dict__

    def __ne__(self, other):
        """Returns true if both objects are not equal"""
        return not self == other
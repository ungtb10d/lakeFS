"""
    lakeFS API

    lakeFS HTTP API  # noqa: E501

    The version of the OpenAPI document: 0.1.0
    Contact: services@treeverse.io
    Generated by: https://openapi-generator.tech
"""


import unittest

import lakefs_client
from lakefs_client.api.templates_api import TemplatesApi  # noqa: E501


class TestTemplatesApi(unittest.TestCase):
    """TemplatesApi unit test stubs"""

    def setUp(self):
        self.api = TemplatesApi()  # noqa: E501

    def tearDown(self):
        pass

    def test_expand_template(self):
        """Test case for expand_template

        """
        pass


if __name__ == '__main__':
    unittest.main()

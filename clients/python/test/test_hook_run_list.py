"""
    lakeFS API

    lakeFS HTTP API  # noqa: E501

    The version of the OpenAPI document: 0.1.0
    Contact: services@treeverse.io
    Generated by: https://openapi-generator.tech
"""


import sys
import unittest

import lakefs_client
from lakefs_client.model.hook_run import HookRun
from lakefs_client.model.pagination import Pagination
globals()['HookRun'] = HookRun
globals()['Pagination'] = Pagination
from lakefs_client.model.hook_run_list import HookRunList


class TestHookRunList(unittest.TestCase):
    """HookRunList unit test stubs"""

    def setUp(self):
        pass

    def tearDown(self):
        pass

    def testHookRunList(self):
        """Test HookRunList"""
        # FIXME: construct object with mandatory attributes with example values
        # model = HookRunList()  # noqa: E501
        pass


if __name__ == '__main__':
    unittest.main()
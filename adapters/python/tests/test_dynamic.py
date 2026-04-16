import pytest
from heimdall.exceptions import HeimdallBaseError
from heimdall.context import hint

class InventorySyncError(HeimdallBaseError):
    category = "INVENTORY_SYNC_FAILED"

def test_heimdall_base_error_category():
    err = InventorySyncError("Failed to sync")
    assert getattr(err, "category", None) == "INVENTORY_SYNC_FAILED"

def test_context_manager_hint():
    try:
        with hint("GOWAY_API_FAIL"):
            raise ValueError("Something went wrong")
    except Exception as e:
        assert getattr(e, "heimdall_category", None) == "GOWAY_API_FAIL"

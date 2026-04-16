from contextlib import contextmanager

@contextmanager
def hint(category: str):
    """Context manager to attach a category to any raised exception (Method B)."""
    try:
        yield
    except Exception as e:
        e.heimdall_category = category
        raise e

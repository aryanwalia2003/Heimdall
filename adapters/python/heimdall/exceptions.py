class HeimdallBaseError(Exception):
    """Base exception for dynamic categorization (Method A).
    Subclasses should override the category attribute.
    """
    category = "SYSTEM_ERROR"

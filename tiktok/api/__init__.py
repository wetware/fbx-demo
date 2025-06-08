import capnp
import logging
import os

logging.basicConfig(level=logging.DEBUG)  # Add this line
logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

_file_path = os.path.abspath(__file__)
_directory = os.path.dirname(_file_path)
capnp_file = os.path.join(_directory, "cap", "tiktok.capnp")

# capnp.remove_import_hook()
tiktok_capnp = capnp.load(capnp_file)
logger.debug("Parsed TikTok capability file.")

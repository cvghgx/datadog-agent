import contextlib
import glob
import json
import os
import platform
import re
import shutil
import string
import sys
import tarfile
import tempfile
from pathlib import Path
from subprocess import check_output

import requests
from invoke import task
from invoke.exceptions import Exit

is_windows = sys.platform == "win32"
is_darwin = sys.platform == "darwin"

from tasks.agent import BUNDLED_AGENTS
from tasks.rtloader import get_dev_path

@task
def build_sds_library(ctx, branch="main"):
    if is_windows:
        printf("not supported")
        return
    with tempfile.TemporaryDirectory() as temp_dir:
        with ctx.cd(temp_dir):
            ctx.run(f"git clone https://github.com/DataDog/dd-sensitive-data-scanner")
            # TODO(remy): checkout a given version
            with ctx.cd("dd-sensitive-data-scanner/sds-go/rust"):
                ctx.run(f"cargo build --release")
                # TODO(remy): add windows support
                dev_path = get_dev_path()
                lib_path = os.path.join(dev_path, "lib")
                if is_darwin:
                    ctx.run(f"cp target/release/libsds_go.dylib {lib_path}")
                else:
                    ctx.run(f"cp target/release/libsds_go.so {lib_path}")

# TODO(remy): clean step


#!/bin/bash
# setup_dev.sh: Setup Heimdall local development environment

set -e

echo "🛡️ Setting up Heimdall Developer Environment..."

# 1. Build the core shared library
echo "🏗️ Building core/libheimdall.so..."
make build

# 2. Setup Python environment
echo "🐍 Setting up Python adapter..."
if [ ! -d "adapters/python/venv" ]; then
    python3 -m venv adapters/python/venv
fi

source adapters/python/venv/bin/activate
pip install -e ./adapters/python

# 3. Create pilot directory
mkdir -p pilot

echo "✅ Environment setup complete!"
echo "🚀 To run the pilot application:"
echo "   export HEIMDALL_LIB_PATH=$(pwd)/bin/libheimdall.so"
echo "   python pilot/main.py"

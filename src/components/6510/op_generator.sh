#!/bin/bash

SRC_DIR="."

if [ ! -d "$SRC_DIR" ]; then
  echo "error: '$SRC_DIR' doesn't exist." >&2
  exit 1
fi

echo
echo "package mos6510"
echo
echo "import \"reflect\""
echo
echo "var _opContainer = map[reflect.Value]string{"
find "$SRC_DIR" -name "inst_*.go" | while read -r file; do
  while IFS= read -r line; do
    func_name=$(echo "${line}" | sed -En 's/^func[ ]+(inst[a-zA-Z0-9]+).+/\1/p')
    [ -z "${func_name}" ] && continue
    echo "    reflect.ValueOf(${func_name}): \"${func_name}\","
  done < "${file}"
done
echo "}"

echo ""
echo "var _opReverseContainer map[string]reflect.Value"
echo ""
echo "func init() {"
echo "    _opReverseContainer = make(map[string]reflect.Value)"
echo "    for f, name := range _opContainer {"
echo "        _opReverseContainer[name] = f"
echo "    }"
echo "}"
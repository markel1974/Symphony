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

echo ""
echo "var _opContainer = make(map[reflect.Value]string)"
echo ""
echo "var _opReverseContainer = make(map[string]reflect.Value)"

echo ""
echo "func init() {"
find "$SRC_DIR" -name "inst_*.go" | while read -r file; do
  while IFS= read -r line; do
    func_name=$(echo "${line}" | sed -En 's/^func[ ]+(inst[a-zA-Z0-9]+).+/\1/p')
    [ -z "${func_name}" ] && continue
    echo "  addOpId(${func_name}, \"${func_name}\")"
  done < "${file}"
done
echo "}"

echo ""
echo "func addOpId(v func(cpu *CPU), id string) {"
echo "	r := reflect.ValueOf(v)"
echo "	_opContainer[r] = id"
echo "	_opReverseContainer[id] = r"
echo "}"

echo ""
echo "func GetOpId(v func(cpu *CPU)) (string, bool) {"
echo "	if v == nil {"
echo "		return "", false"
echo "	}"
echo "	r := reflect.ValueOf(v)"
echo "	ret, ok := _opContainer[r]"
echo "	if !ok {"
echo "		return "", false"
echo "	}"
echo "	return ret, true"
echo "}"

echo ""
echo "func GetOpFn(v string) (func(cpu *CPU), bool) {"
echo "	r, ok := _opReverseContainer[v]"
echo "	if !ok {"
echo "		return nil, false"
echo "	}"
echo "	ret, ok := r.Interface().(func(cpu *CPU))"
echo "	if !ok {"
echo "		return nil, false"
echo "	}"
echo "	return ret, true"
echo "}"
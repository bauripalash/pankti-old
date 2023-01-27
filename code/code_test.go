package code

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       OpCode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpAdd, []int{}, []byte{byte(OpAdd)}},
	}

	expects := []string{
		"0000 OpConstant 65534\n",
		"0000 OpAdd\n",
	}

	for idx, tt := range tests {
		instrs := Make(tt.op, tt.operands...)
		if len(instrs) != len(tt.expected) {
			t.Errorf("Instruction with wrong length; W=%d , G=%d", len(tt.expected), len(instrs))
		}

		for i, b := range tt.expected {
			if instrs[i] != tt.expected[i] {
				t.Errorf("Wrong Byte as pos %d; W=%d G=%d", i, b, instrs[i])
			}
		}

		str_expected := expects[idx] //"0000 OpConstant 65534\n"
		ins := Instructions(Make(tt.op, tt.operands...))

		if str_expected != ins.String() {
			t.Errorf("String Format -> Wanted = %s\nGot = %s", str_expected, ins)
		}

	}

}

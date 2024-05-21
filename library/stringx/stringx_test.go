package stringx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSnakeString(t *testing.T) {
	data := [][2]string{
		{"XxYy", "xx_yy"},
		{"_XxYy", "_xx_yy"},
		{"TcpRpc", "tcp_rpc"},
		{"ID", "id"},
		{"UserID", "user_id"},
		{"RPC", "rpc"},
		{"TCP_RPC", "tcp_rpc"},
		{"wakeRPC", "wake_rpc"},
		{"_TCP__RPC", "_tcp__rpc"},
		{"_TcP__RpC_", "_tc_p__rp_c_"},
	}
	for _, p := range data {
		r := SnakeString(p[0])
		assert.Equal(t, p[1], r, p[0])
		r = SnakeString(p[1])
		assert.Equal(t, p[1], r, p[0])
	}
}

func TestCamelString(t *testing.T) {
	data := [][2]string{
		{"_", "_"},
		{"xx_yy", "XxYy"},
		{"_xx_yy", "_XxYy"},
		{"id", "Id"},
		{"user_id", "UserId"},
		{"rpc", "Rpc"},
		{"tcp_rpc", "TcpRpc"},
		{"wake_rpc", "WakeRpc"},
		{"_tcp___rpc", "_Tcp__Rpc"},
		{"_tc_p__rp_c__", "_TcP_RpC__"},
	}
	for _, p := range data {
		r := CamelString(p[0])
		assert.Equal(t, p[1], r, p[0])
		r = CamelString(p[1])
		assert.Equal(t, p[1], r, p[0])
	}
}

func TestLintCamelString(t *testing.T) {
	data := [][2]string{
		{"_", "_"},
		{"xx_yy", "XxYy"},
		{"_xx_yy", "XxYy"},
		{"id", "ID"},
		{"user_id", "UserID"},
		{"rpc", "RPC"},
		{"tcp_rpc", "TCPRPC"},
		{"wake_rpc", "WakeRPC"},
		{"___tcp___rpc", "TCPRPC"},
		{"_tc_p__rp_c__", "TcPRpC"},
	}
	for _, p := range data {
		r := LintCamelString(p[0])
		assert.Equal(t, p[1], r, p[0])
		r = LintCamelString(p[1])
		assert.Equal(t, p[1], r, p[0])
	}
}

func TestHTMLEntityToUTF8(t *testing.T) {
	want := `{"info":[["color","咖啡色|绿色"]]｝`
	got := HTMLEntityToUTF8(`{"info":[["color","&#5496;&#5561;&#8272;&#7c;&#7eff;&#8272;"]]｝`, 16)
	if got != want {
		t.Fatalf("want: %q, got: %q", want, got)
	}
}

func TestCodePointToUTF8(t *testing.T) {
	got := CodePointToUTF8(`{"info":[["color","\u5496\u5561\u8272\u7c\u7eff\u8272"]]｝`, 16)
	want := `{"info":[["color","咖啡色|绿色"]]｝`
	if got != want {
		t.Fatalf("want: %q, got: %q", want, got)
	}
}

func TestSpaceInOne(t *testing.T) {
	a := struct {
		input  string
		output string
	}{
		input: `# authenticate method 

		//  comment2	

		/*  some other 
			  comments */
		`,
		output: `# authenticate method
	// comment2
	/* some other
	comments */
	`,
	}
	r := SpaceInOne(a.input)
	if r != a.output {
		t.Fatalf("want: %q, got: %q", a.output, r)
	}
}

func TestExampleStringMarshalJSON(t *testing.T) {
	s := `<>&{}""`
	t.Logf("%s\n", StringMarshalJSON(s, true))
	t.Logf("%s\n", StringMarshalJSON(s, false))
	// Output:
	// "\u003c\u003e\u0026{}\"\""
	// "<>&{}\"\""
}

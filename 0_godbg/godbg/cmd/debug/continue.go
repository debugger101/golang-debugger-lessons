package debug

import (
	"fmt"
	"syscall"

	"github.com/spf13/cobra"
)

var continueCmd = &cobra.Command{
	Use:   "continue",
	Short: "运行到下个断点",
	Annotations: map[string]string{
		cmdGroupAnnotation: cmdGroupCtrlFlow,
	},
	Aliases: []string{"c"},
	RunE: func(cmd *cobra.Command, args []string) error {
		//fmt.Println("continue")
		// 读取PC值
		regs := syscall.PtraceRegs{}
		err := syscall.PtraceGetRegs(TraceePID, &regs)
		if err != nil {
			return fmt.Errorf("get regs error: %v", err)
		}

		buf := make([]byte, 1)
		n, err := syscall.PtracePeekText(TraceePID, uintptr(regs.PC()), buf)
		if err != nil || n != 1 {
			return fmt.Errorf("peek text error: %v, bytes: %d", err, n)
		}

		// read a breakpoint
		if buf[0] == 0xCC {
			regs.SetPC(regs.PC() - 1)
			// TODO refactor breakpoint.Disable()/Enable() methods
			orig := breakpoints[uintptr(regs.PC())].Orig
			n, err := syscall.PtracePokeText(TraceePID, uintptr(regs.PC()), []byte{orig})
			if err != nil || n != 1 {
				return fmt.Errorf("poke text error: %v, bytes: %d", err, n)
			}
		}

		err = syscall.PtraceCont(TraceePID, 0)
		if err != nil {
			return fmt.Errorf("single step error: %v", err)
		}

		// MUST: 当发起了某些对tracee执行控制的ptrace request之后，要调用syscall.Wait等待并获取tracee状态变化
		var wstatus syscall.WaitStatus
		var rusage syscall.Rusage
		_, err = syscall.Wait4(TraceePID, &wstatus, syscall.WALL, &rusage)
		if err != nil {
			return fmt.Errorf("wait error: %v", err)
		}

		// display current pc
		regs = syscall.PtraceRegs{}
		err = syscall.PtraceGetRegs(TraceePID, &regs)
		if err != nil {
			return fmt.Errorf("get regs error: %v", err)
		}
		fmt.Printf("continue ok, current PC: %#x\n", regs.PC())
		return nil
	},
}

func init() {
	debugRootCmd.AddCommand(continueCmd)
}

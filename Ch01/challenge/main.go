package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

/*
Instructions:

- Pass text to print as an argument
  - If not argument - read from stdin

- Use -width to specify width
  - Width should be bigger than 0 and less than 250
  - Default to 80

- Use -out to specify output file
  - Default to stdout
*/

type Width struct {
  val *int
}

func (w *Width) String() string {
  if w.val == nil {
    return ""
  }

  return fmt.Sprintf("%d", *w.val)
}

func (w *Width) Set(v string) error {
  n,err := strconv.Atoi(v)
  if err != nil {
    return fmt.Errorf("bad number: %s", err)
  }

  if n<=0 || n>= 250{
    return fmt.Errorf("width must be between 0 and 250")
  }

  *w.val=n

  return nil
}

type out struct {
  val *io.Writer
  file *os.File
}

func (o *out) String() string {
  if o.val == nil {
    return ""
  }
  if o.file != nil {
		return o.file.Name()
	}

  return fmt.Sprintf("%d", *o.val)
}

func (o *out) Set(v string) error {
  if v== ""{
    *o.val=os.Stdout
    return nil
  }
  file,err:=os.OpenFile(v,os.O_RDWR|os.O_CREATE, 0644)
  if err !=nil{
    return fmt.Errorf("could not open file %s error:%s",v,err)
  }

  *o.val=file
  o.file=file

  return nil
}

func (o *out) Close()  {
  if o.file != nil{
    o.file.Close()
  }
}

var config struct {
  width int
  out  io.Writer
}
func main() {
  config.width=80
  config.out=os.Stdout
  outFlag := &out{val: &config.out}
  flag.Var(&Width{&config.width},"width","width of banner default 80")
  flag.Var(outFlag,"output", "output file to print to")
  flag.Parse()

  defer outFlag.Close()

  text:=""
  switch flag.NArg() {
	case 0:
		d,err:=io.ReadAll(os.Stdin)
    if err !=nil {
      fmt.Fprintf(os.Stderr, "error: can't read - %s\n", err)
			os.Exit(1)
    }
    text=string(d)
	case 1:
		text = flag.Arg(0)
	default:
		fmt.Fprintln(os.Stderr, "error: wrong number of arguments")
		os.Exit(1)
	}
	Banner(config.out, text, config.width)
}

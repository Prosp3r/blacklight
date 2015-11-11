package main

import (
	"fmt"
	"strconv"
)

type operation interface {
	Eval(stack) stack
	Value() datatypes
	String() string
}

type Op struct {
	Name string
	Data datatypes
}

func (o Op) Eval(current stack) stack {
	switch o.Name {

	// @stack (current Stack)
	case "decap":
		current.Decap()
	case "depth":
		current.Push(NewNumber(current.Depth()))
	case "drop":
		current.Drop()
	case "dup":
		current.Dup()
	case "over":
		current.Over()
	case "purge":
		current.Purge()
	case "rot":
		current.Rot()
	case "swap":
		current.Swap()

	// NativeIntegers (Number)
	case "add":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		sum := n2 + n1
		current.Push(NewNumber(sum))
	case "sub":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		result := n2 - n1
		current.Push(NewNumber(result))
	case "mul":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		product := n2 * n1
		current.Push(NewNumber(product))
	case "div":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		result := n2 / n1
		current.Push(NewNumber(result))
	case "mod":
		i1 := current.Pop()
		i2 := current.Pop()
		n1 := i1.Value().(int)
		n2 := i2.Value().(int)
		remainder := n2 % n1
		current.Push(NewNumber(remainder))
	case "n-to-c":
		i := current.Pop()
		n := i.Value().(int)
		str := strconv.Itoa(n)
		current.Push(NewChar("\\" + str))
	case "n-to-cv":
		i := current.Pop()
		n := i.Value().(int)
		str := strconv.Itoa(n)
		current.Push(NewCharVector(str))

	// Chars
	case "c-to-cv":
		c := current.Pop()
		current.Push(NewCharVector(c.(Char).C_to_T()))
	case "c-to-n":
		c := current.Pop().(Char)
		current.Push(NewNumber(c.C_to_N()))

	// Debug
	case "print":
		i := current.Pop()
		fmt.Println(i.String())

	// Vectors
	case "()":
		fallthrough
	case "v-new":
		current.Push(NewVector([]datatypes{}))
	case "''":
		current.Push(NewCharVector("''"))
	case "cat":
		i1 := current.Pop().(vector)
		i2 := current.Pop().(vector)
		result := i2.Cat(i1)
		current.Push(result)
	case "app":
		i := current.Pop()
		v := current.Pop().(vector)
		current.Push(v.App(i))
	case "ato":
		n := current.Pop().(*Number)
		v := current.Peek().(vector)
		i := v.Ato(n.Value().(int))
		current.Push(i)
	case "rmo":
		n := current.Pop().(*Number).Value().(int)
		v := current.Pop().(vector)
		nv := v.Rmo(n)
		current.Push(nv)
	case "len":
		v := current.Peek().(vector)
		current.Push(NewNumber(v.Len()))

	// Queues
	case "newq":
		q := NewQueue()
		current.Push(q)
	case "enq":
		i := current.Pop()
		q := current.Peek().(*Queue)
		q.Enqueue(i)
	case "deq":
		q := current.Peek().(*Queue)
		i := q.Dequeue()
		current.Push(i)
	case "q-to-v":
		q_to_v(current.(*Stack))
	case "q-to-cv":
		q_to_cv(current.(*Stack))

	// Stacks
	case "news":
		fallthrough
	case "<>":
		fallthrough
	case "s-new":
		s := NewStack("user")
		current.Push(s)
	case "push":
		i := current.Pop()
		s := current.Peek().(stack)
		s.Push(i)
	case "pop":
		s := current.Peek().(stack)
		current.Push(s.Pop())
	case "size":
		s := current.Peek().(stack)
		current.Push(NewNumber(s.Depth()))
	case "tail":
		s := current.Peek().(stack)
		s.Drop()

	// Logic
	case "eq":
		i1 := current.Pop()
		i2 := current.Peek()
		current.Push(blEq(i1, i2))
	case "not":
		t := current.Pop()
		switch t.(type) {
		case *Tag:
			if t.(*Tag).Kind == "nil" {
				current.Push(NewTrue("not"))
			} else {
				current.Push(NewNil("not"))
			}
		default:
			current.Push(NewNil("not"))
		}

	case "nil":
		current.Push(NewNil("nil"))
	case "true":
		current.Push(NewTrue("true"))

	// IO
	case "read":
		source := current.Pop()
		q := current.Peek().(*Queue)
		io := ReadIO(source, q)
		current.Push(io)
	case "write":
		dest := current.Pop()
		q := current.Peek().(*Queue)
		io := WriteIO(dest, q)
		current.Push(io)

	// Objects
	case "o-new":
		current.Push(NewObject())
	case "set":
		slot := current.Pop().(Word)
		i := current.Pop()
		o := current.Peek().(*Object)
		o.Set(slot, i)
	case "fetch":
		slot := current.Pop().(Word)
		o := current.Peek().(*Object)
		current.Push(o.Fetch(slot))
	case "child":
		o := current.Peek().(*Object)
		child := NewChildObject(o)
		current.Push(child)

	default:
		warn("UNIMPLEMENTED operation: " + o.String())
	}
	return current
}

func (o Op) Value() datatypes {
	return o.Data
}

func (o Op) String() string {
	return o.Name
}

func newOp(t string) *Op {
	op := new(Op)
	op.Name = t
	return op
}

type metaOp struct {
	Op
}

func (m metaOp) Eval(meta stack) stack {
	switch m.Name {
	case "@":
		s := meta.Peek()
		current := s.(*Stack)
		current.Push(current)
	case "^":
		s := meta.Peek()
		current := s.(*Stack)
		meta.Swap()
		s = meta.Peek()
		prev := s.(*Stack)
		meta.Swap()
		current.Push(prev)
	case "$":
		s := meta.Peek()
		current := s.(*Stack)
		current.Push(meta)
	case "$decap":
		meta.Decap()
	case "$drop":
		meta.Drop()
	case "$new":
		if meta.Depth() > 0 {
			s := meta.Peek()
			os := s.(*Stack)
			ns := NewSystemStack()
			ns.Push(os)
			meta.Push(ns)
		} else {
			meta.Push(NewSystemStack())
			meta = newMetaOp("$new").Eval(meta)
		}
	case "$swap":
		meta.Swap()

	// Loops and Logic
	case "until":
		current := meta.Peek().(*Stack)
		comp := current.Pop().(WordVector).Ops
		actn := current.Pop().(WordVector).Ops
	Until:
		for {
			doEval(meta.(*MetaStack), comp)
			current = meta.Peek().(*Stack)
			if current.Pop().(*Tag).Kind == "true" {
				break Until
			}
			doEval(meta.(*MetaStack), actn)
		}
	case "loop":
		current := meta.Peek().(*Stack)
		actn := current.Pop().(WordVector).Ops
		for {
			doEval(meta.(*MetaStack), actn)
		}
	case "proq":
		processQueue(meta.(*MetaStack))
	case "if":
		current := meta.Peek().(*Stack)
		comp := current.Pop().(WordVector).Ops
		actn := current.Pop().(WordVector).Ops
		doEval(meta.(*MetaStack), comp)
		current = meta.Peek().(*Stack)
		if current.Pop().(*Tag).Kind == "true" {
			doEval(meta.(*MetaStack), actn)
		}
	case "either":
		current := meta.Peek().(*Stack)
		comp := current.Pop().(WordVector).Ops
		iffalse := current.Pop().(WordVector).Ops
		iftrue := current.Pop().(WordVector).Ops
		doEval(meta.(*MetaStack), comp)
		current = meta.Peek().(*Stack)
		if current.Pop().(*Tag).Kind == "true" {
			doEval(meta.(*MetaStack), iftrue)
		} else {
			doEval(meta.(*MetaStack), iffalse)
		}

	// File Loading & Multithreading
	case "do":
		current := meta.Peek().(*Stack)
		filename := current.Pop().(*CharVector).Value().(string)
		code := loadFile(filename)
		tokens := parse(code)
		ops := lex(tokens)
		doEval(meta.(*MetaStack), ops)
	case "co":
		co(meta.(*MetaStack))
	case "wait":
		threads.Wait()
	case "bkg":
		bkg(meta.(*MetaStack))
	case "work":
		work(meta.(*MetaStack))

	// Objects & Eval
	case "get":
		// Disabled due to Meta/MetaStack confusion
		//current := meta.Peek().(*Stack)
		//slot := current.Pop().(Word)
		//o := current.Peek().(*Object)
		//o.Get(meta.(*MetaStack), slot)
	case "self":
		current := meta.Peek().(*Stack)
		o := *meta.(*MetaStack).ObjectStack.Peek()
		current.Push(o)
	case "call":
		current := meta.Peek().(*Stack)
		wv := current.Pop().(WordVector)
		wv.Call(meta.(*MetaStack))

	default:
		warn("UNIMPLEMENTED $operation: " + m.String())
	}

	return meta
}

func newMetaOp(t string) *metaOp {
	op := new(metaOp)
	op.Name = t
	return op
}

type pushLiteral struct {
	Op
}

func (o pushLiteral) Eval(s stack) stack {
	s.Push(o.Value())
	return s
}

func (o pushLiteral) String() string {
	return fmt.Sprint(o.Value())
}

type pushNumber struct {
	pushLiteral
}

func newPushNumber(t string) *pushNumber {
	pi := new(pushNumber)
	pi.Name = t
	i, _ := strconv.Atoi(t)
	pi.Data = NewNumber(i)
	return pi
}

type pushWord struct {
	pushLiteral
}

func newPushWord(t string) *pushWord {
	var w Word
	pw := new(pushWord)

	if t[0] == ":"[0] || t[0] == "~"[0] {
		w = NewWord(t[1:])
	} else if t[len(t)-1] == ":"[0] {
		w = NewWord(t[:len(t)-1])
	} else {
		w = NewWord(t)
	}

	pw.Name = t
	pw.Data = w
	return pw
}

type pushVector struct {
	pushLiteral
	Contents []operation
}

func newPushVector(t string) *pushVector {
	pv := new(pushVector)
	pv.Name = t
	return pv
}

func (pv *pushVector) Eval(s stack) stack {
	var data []datatypes
	for _, op := range pv.Contents {
		data = append(data, op.Value())
	}
	v := NewVector(data)
	s.Push(v)
	return s
}

type pushWordVector struct {
	pushVector
}

func newPushWordVector(t string) *pushWordVector {
	pwv := new(pushWordVector)
	pwv.Name = t
	return pwv
}

func (pwv *pushWordVector) Eval(s stack) stack {
	wv := NewWordVector(pwv.Contents)
	s.Push(wv)
	return s
}

type pushCharVector struct {
	pushLiteral
	Contents []operation
}

func newPushCharVector(t string) *pushCharVector {
	ps := new(pushCharVector)
	ps.Name = t
	ps.Data = NewCharVector(t)
	return ps
}

type pushChar struct {
	pushLiteral
}

func newPushChar(t string) *pushChar {
	pc := new(pushChar)
	pc.Name = t
	pc.Data = NewChar(t)
	return pc
}

type pushQueue struct {
	pushLiteral
	Contents []operation
}

func processQueue(meta *MetaStack) {
	current := meta.Peek().(*Stack)
	wv := current.Pop().(WordVector)
	q := current.Pop().(*Queue)
	var tokens []string

	for _, w := range wv.Data {
		tokens = append(tokens, w.Name)
	}

ProcLoop:
	for {
		select {
		case item := <-q.Items:
			current.Push(item)
			ops := lex(tokens)
			doEval(meta, ops)
		default:
			break ProcLoop
		}
	}
}

func q_to_v(s *Stack) {
	q := s.Pop().(*Queue)
	items := []datatypes{}

QVLoop:
	for {
		select {
		case i := <-q.Items:
			items = append(items, i)
		default:
			break QVLoop
		}
	}

	v := NewVector(items)
	s.Push(v)
}

func q_to_cv(s *Stack) {
	q := s.Pop().(*Queue)
	str := ""

QVLoop:
	for {
		i := <-q.Items
		if blEq(i, NewNil("q_to_cv")).Bool() {
			break QVLoop
		} else {
			str = str + i.(Char).CVString()
		}
	}

	v := NewCharVector(str)
	s.Push(v)
}

func bkg(meta *MetaStack) {
	current := meta.Peek().(*Stack)
	wv := current.Pop().(WordVector)
	i := current.Pop()

	threads.Add(1)
	go func(item datatypes) {
		defer threads.Done()
		new_meta := NewMetaStack()
		new_current := NewStack("work")
		new_meta.Push(new_current)
		new_current.Push(item)
		doEval(new_meta, wv.Ops)
	}(i)
}

func work(meta *MetaStack) {
	current := meta.Peek().(*Stack)
	wv := current.Pop().(WordVector)
	in := current.Pop().(*Queue)
	out := current.Pop().(*Queue)

	threads.Add(1)
	go func(in *Queue, out *Queue) {
		defer threads.Done()
		new_meta := NewMetaStack()
		new_current := NewStack("work")
		new_meta.Push(new_current)
		new_current.Push(in)
		new_current.Push(out)
		doEval(new_meta, wv.Ops)
	}(in, out)

	current.Push(out)
	current.Push(in)
}

func co(meta *MetaStack) {
	current := meta.Peek().(*Stack)
	filename := current.Pop().(*CharVector).Value().(string)
	in := NewQueue()
	out := NewQueue()

	threads.Add(1)
	go func(filename string, in *Queue, out *Queue) {
		defer threads.Done()
		code := loadFile(filename)
		tokens := parse(code)
		ops := lex(tokens)
		new_meta := NewMetaStack()
		new_current := NewStack("co")
		new_meta.Push(new_current)
		new_current.Push(in)
		new_current.Push(out)
		doEval(new_meta, ops)
	}(filename, in, out)

	current.Push(out)
	current.Push(in)
}

package main

import (
	"fmt"
	"unsafe"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
)

var (
	// Keep track of debug locations
	dbgLoc = make(map[string]*metadata.DILocation)
)

func main() {
	// Convenience constants.
	var (
		emptyExpr = &metadata.DIExpression{
			MetadataID: -1,
		}
		i8nullptr = constant.Null{Typ: types.I8Ptr}
		null      = constant.NewNull(&types.PointerType{})
		one       = constant.NewInt(types.I32, 1)
		ten       = constant.NewInt(types.I8, 10)
		twelve    = constant.NewInt(types.I32, 12)
		thirteen  = constant.NewInt(types.I8, 13)
		zero      = constant.NewInt(types.I32, 0)
	)

	// !DIExpression()

	// Create LLVM IR module.
	m := ir.NewModule()

	// Add metadata
	addMetadata(m)

	//  ; ModuleID = 'main.go'
	//  source_filename = "main.go"
	//  target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
	//  target triple = "wasm32-unknown-unknown-wasm"
	m.DataLayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
	m.TargetTriple = "wasm32-unknown-unknown-wasm"
	m.SourceFilename = "main.go"

	// 	attributes #0 = { optsize }
	atGrp0 := ir.AttrGroupDef{ID: 0, FuncAttrs: []ir.FuncAttribute{
		enum.FuncAttrOptSize},
	}
	m.AttrGroupDefs = append(m.AttrGroupDefs, &atGrp0)

	// 	attributes #1 = { noreturn optsize }
	atGrp1 := ir.AttrGroupDef{ID: 1, FuncAttrs: []ir.FuncAttribute{
		enum.FuncAttrNoReturn,
		enum.FuncAttrOptSize},
	}
	m.AttrGroupDefs = append(m.AttrGroupDefs, &atGrp1)

	// 	attributes #2 = { cold noreturn nounwind optsize }
	atGrp2 := ir.AttrGroupDef{ID: 2, FuncAttrs: []ir.FuncAttribute{
		enum.FuncAttrCold,
		enum.FuncAttrNoReturn,
		enum.FuncAttrNoUnwind,
		enum.FuncAttrOptSize},
	}
	m.AttrGroupDefs = append(m.AttrGroupDefs, &atGrp2)

	// 	attributes #3 = { nounwind optsize readnone speculatable }
	atGrp3 := ir.AttrGroupDef{ID: 3, FuncAttrs: []ir.FuncAttribute{
		enum.FuncAttrNoUnwind,
		enum.FuncAttrOptSize,
		enum.FuncAttrReadNone,
		enum.FuncAttrSpeculatable},
	}
	m.AttrGroupDefs = append(m.AttrGroupDefs, &atGrp3)

	// 	attributes #4 = { argmemonly nounwind optsize readonly }
	atGrp4 := ir.AttrGroupDef{ID: 4, FuncAttrs: []ir.FuncAttribute{
		enum.FuncAttrArgMemOnly,
		enum.FuncAttrNoUnwind,
		enum.FuncAttrOptSize,
		enum.FuncAttrReadOnly},
	}
	m.AttrGroupDefs = append(m.AttrGroupDefs, &atGrp4)

	// 	attributes #5 = { nounwind }
	atGrp5 := ir.AttrGroupDef{ID: 5, FuncAttrs: []ir.FuncAttribute{
		enum.FuncAttrNoUnwind},
	}
	m.AttrGroupDefs = append(m.AttrGroupDefs, &atGrp5)

	// * Globals *

	// 	@runtime.runqueueBack = internal unnamed_addr global i8* null
	globalRunQueueBack := m.NewGlobalDef("runtime.runqueueBack", &i8nullptr)
	globalRunQueueBack.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalRunQueueBack.Linkage = enum.LinkageInternal

	// 	@runtime.runqueueFront = internal unnamed_addr global i8* null
	globalRunQueueFront := m.NewGlobalDef("runtime.runqueueFront", &i8nullptr)
	globalRunQueueFront.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalRunQueueFront.Linkage = enum.LinkageInternal

	// 	@runtime.stdout = internal unnamed_addr global i32 0
	globalStdout := m.NewGlobalDef("runtime.stdout", zero)
	globalStdout.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalStdout.Linkage = enum.LinkageInternal

	// 	@"runtime.nilPanic$string" = internal unnamed_addr constant [23 x i8] c"nil pointer dereference"
	globalNilPanic := m.NewGlobalDef("runtime.nilPanic$string", constant.NewCharArrayFromString("nil pointer dereference"))
	globalNilPanic.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalNilPanic.Linkage = enum.LinkageInternal

	// 	@"runtime.runtimePanic$string" = internal unnamed_addr constant [22 x i8] c"panic: runtime error: "
	globalRuntimePanic := m.NewGlobalDef("runtime.runtimePanic$string", constant.NewCharArrayFromString("panic: runtime error: "))
	globalRuntimePanic.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalRuntimePanic.Linkage = enum.LinkageInternal

	// Define some needed function parameters early
	activateTaskParam0 := ir.NewParam("", types.I8Ptr)
	activateTaskParamContext := ir.NewParam("context", types.I8Ptr)
	activateTaskParamParHand := ir.NewParam("parentHandle", types.I8Ptr)
	coroSubfnAddrParam0 := ir.NewParam("", types.I8Ptr)
	coroSubfnAddrParam1 := ir.NewParam("", types.I8)
	dbgValParam0 := ir.NewParam("", types.Metadata)
	dbgValParam1 := ir.NewParam("", types.Metadata)
	dbgValParam2 := ir.NewParam("", types.Metadata)
	memSetParam0 := ir.NewParam("", types.I8Ptr)
	memSetParam1 := ir.NewParam("", types.I8)
	memSetParam2 := ir.NewParam("", types.I32)
	prtStrParam0 := ir.NewParam("", types.I8Ptr)
	prtStrParam1 := ir.NewParam("", types.I32)
	putCharParam0 := ir.NewParam("", types.I8)
	resWriParam0 := ir.NewParam("", types.I32)
	resWriParam1 := ir.NewParam("", types.I8Ptr)
	resWriParam2 := ir.NewParam("", types.I32)

	// Add functions here, so they're emitted in the same order as the target ll
	startFunc := m.NewFunc("_start", types.Void)
	activateTaskFunc := m.NewFunc("runtime.activateTask", types.Void, activateTaskParam0, activateTaskParamContext, activateTaskParamParHand)
	cwaMainFunc := m.NewFunc("cwa_main", types.Void)
	getFuncPtrFunc := m.NewFunc("runtime.getFuncPtr", types.Void)
	IOGetStdoutFunc := m.NewFunc("io_get_stdout", types.I32)
	memSetFunc := m.NewFunc("memset", types.I8Ptr, memSetParam0, memSetParam1, memSetParam2)
	nilPanicFunc := m.NewFunc("runtime.nilPanic", types.Void)
	prtNlFunc := m.NewFunc("runtime.printnl", types.Void)
	prtStrFunc := m.NewFunc("runtime.printstring", types.Void, prtStrParam0, prtStrParam1)
	putCharFunc := m.NewFunc("runtime.putchar", types.Void, putCharParam0)
	resWriFunc := m.NewFunc("resource_write", types.I32, resWriParam0, resWriParam1, resWriParam2)
	resumeFunc := m.NewFunc("resume", types.Void)
	runtimePanicFunc := m.NewFunc("runtime.runtimePanic", types.Void)
	LLVMTrapFunc := m.NewFunc("llvm.trap", types.Void)
	LLVMDbgValueFunc := m.NewFunc("llvm.dbg.value", types.Void, dbgValParam0, dbgValParam1, dbgValParam2)
	LLVMCoroSubFnAddrFunc := m.NewFunc("llvm.coro.subfn.addr", types.I8Ptr, coroSubfnAddrParam0, coroSubfnAddrParam1)

	// 	; Function Attrs: optsize
	// 	define void @_start() local_unnamed_addr #0 section ".text._start" !dbg !5 {
	startFunc.UnnamedAddr = enum.UnnamedAddrLocalUnnamedAddr
	startFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	startFunc.Section = ".text._start"
	startFunc.Metadata = append(startFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["_start"]})

	// 	entry:
	startEntryBlock := startFunc.NewBlock("entry")

	// 	%0 = tail call i32 @io_get_stdout(), !dbg !8
	startVal0 := startEntryBlock.NewCall(IOGetStdoutFunc)
	startVal0.Tail = enum.TailTail
	startVal0.Metadata = append(startVal0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["8"]})

	// 	store i32 %0, i32* @runtime.stdout, align 4, !dbg !8
	startStore0 := startEntryBlock.NewStore(startVal0, globalStdout)
	startStore0.Align = 4
	startStore0.Metadata = append(startStore0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["8"]})

	// 	ret void, !dbg !12
	startRet0 := startEntryBlock.NewRet(nil)
	startRet0.Metadata = append(startRet0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["12"]})
	// 	}

	// 	@"main.go.main$string" = internal unnamed_addr constant [12 x i8] c"Hello world!"
	globalMain := m.NewGlobalDef("main.go.main$string", constant.NewCharArrayFromString("Hello world!"))
	globalMain.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	globalMain.Linkage = enum.LinkageInternal

	// * Function declarations *

	// 	; Function Attrs: optsize
	// 	declare i32 @io_get_stdout() local_unnamed_addr #0
	IOGetStdoutFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	IOGetStdoutFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}

	// 	; Function Attrs: cold noreturn nounwind optsize
	// 	declare void @llvm.trap() #2
	LLVMTrapFunc.FuncAttrs = []ir.FuncAttribute{&atGrp2}

	// 	; Function Attrs: nounwind optsize readnone speculatable
	// 	declare void @llvm.dbg.value(metadata, metadata, metadata) #3
	LLVMDbgValueFunc.FuncAttrs = []ir.FuncAttribute{&atGrp3}

	// 	; Function Attrs: argmemonly nounwind optsize readonly
	// 	declare i8* @llvm.coro.subfn.addr(i8* nocapture readonly, i8) #4
	coroSubfnAddrParam0.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture, enum.ParamAttrReadOnly}
	LLVMCoroSubFnAddrFunc.FuncAttrs = []ir.FuncAttribute{&atGrp4}

	// 	; Function Attrs: optsize
	// 	declare i32 @resource_write(i32, i8* nocapture, i32) local_unnamed_addr #0
	resWriParam1.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture}
	resWriFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	resWriFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}

	// * Function definitions *

	// 	; Function Attrs: optsize
	// 	define dso_local void @runtime.activateTask(i8*, i8* nocapture readnone %context, i8* nocapture readnone %parentHandle) unnamed_addr #0 section ".text.runtime.activateTask" !dbg !13 {
	activateTaskParamContext.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture, enum.ParamAttrReadNone}
	activateTaskParamParHand.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture, enum.ParamAttrReadNone}
	activateTaskFunc.Preemption = enum.PreemptionDSOLocal
	activateTaskFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	activateTaskFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	activateTaskFunc.Section = ".text.runtime.activateTask"
	activateTaskFunc.Metadata = append(activateTaskFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.activateTask"]})

	// Create the blocks for the upcoming branches
	activateTaskEntryBlock := activateTaskFunc.NewBlock("entry")
	activateTaskIfThenBlock := activateTaskFunc.NewBlock("if.then")
	activateTaskIfDoneBlock := activateTaskFunc.NewBlock("if.done")
	activateTaskIfThenIBlock := activateTaskFunc.NewBlock("if.then.i")
	activateTaskIfDone3IBlock := activateTaskFunc.NewBlock("if.done3.i")
	activateTaskIfThen4IBlock := activateTaskFunc.NewBlock("if.then4.i")
	activateTaskStoreNextIBlock := activateTaskFunc.NewBlock("store.next.i")

	// 	entry:
	// 	  call void @llvm.dbg.value(metadata i8* %0, metadata !21, metadata !DIExpression()), !dbg !22
	activateTaskCall0 := activateTaskEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: activateTaskParam0}, &metadata.Value{Value: dbgLoc["LocalVarTask"]}, &metadata.Value{Value: emptyExpr})
	activateTaskCall0.Metadata = append(activateTaskCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["22"]})

	// 	  %1 = icmp eq i8* %0, null, !dbg !23
	activateTaskVal1 := activateTaskEntryBlock.NewICmp(enum.IPredEQ, activateTaskParam0, null)
	activateTaskVal1.Metadata = append(activateTaskVal1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["23"]})

	// 	  br i1 %1, label %if.then, label %if.done, !dbg !24
	activateTaskBr0 := activateTaskEntryBlock.NewCondBr(activateTaskVal1, activateTaskIfThenBlock, activateTaskIfDoneBlock)
	activateTaskBr0.Metadata = append(activateTaskBr0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["24"]})

	//  if.then:                                          ; preds = %if.then.i, %store.next.i, %if.then4.i, %entry
	// 	  ret void, !dbg !25
	activateTaskBr1 := activateTaskIfThenBlock.NewRet(nil)
	activateTaskBr1.Metadata = append(activateTaskBr1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["25"]})

	// 	if.done:                                          ; preds = %entry
	// 	  call void @llvm.dbg.value(metadata i8* %0, metadata !26, metadata !DIExpression()), !dbg !29
	activateTaskCall1 := activateTaskIfDoneBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: activateTaskParam0}, &metadata.Value{Value: dbgLoc["LocalVarT"]}, &metadata.Value{Value: emptyExpr})
	activateTaskCall1.Metadata = append(activateTaskCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["29"]})

	// 	  %2 = bitcast i8* %0 to i8**, !dbg !31
	activateTaskVal2 := activateTaskIfDoneBlock.NewBitCast(activateTaskParam0, types.NewPointer(types.I8Ptr))
	activateTaskVal2.Metadata = append(activateTaskVal2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["31"]})

	// 	  %3 = load i8*, i8** %2, align 4, !dbg !31
	activateTaskVal3 := activateTaskIfDoneBlock.NewLoad(activateTaskVal2)
	activateTaskVal3.Align = 4
	activateTaskVal3.Metadata = append(activateTaskVal3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["31"]})

	// 	  %4 = icmp eq i8* %3, null, !dbg !31
	activateTaskVal4 := activateTaskIfDoneBlock.NewICmp(enum.IPredEQ, activateTaskVal3, null)
	activateTaskVal4.Metadata = append(activateTaskVal4.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["31"]})

	// 	  br i1 %4, label %if.then.i, label %if.done3.i, !dbg !32
	activateTaskBr2 := activateTaskIfDoneBlock.NewCondBr(activateTaskVal4, activateTaskIfThenIBlock, activateTaskIfDone3IBlock)
	activateTaskBr2.Metadata = append(activateTaskBr2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["32"]})

	// 	if.then.i:                                        ; preds = %if.done
	// 	  %5 = bitcast i8* %0 to { i8*, i8* }*
	activateTaskVal5 := activateTaskIfThenIBlock.NewBitCast(activateTaskParam0, types.NewPointer(types.NewStruct(types.I8Ptr, types.I8Ptr)))

	// 	  %6 = getelementptr inbounds { i8*, i8* }, { i8*, i8* }* %5, i32 0, i32 1
	activateTaskVal6 := activateTaskIfThenIBlock.NewGetElementPtr(activateTaskVal5, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 1))
	activateTaskVal6.InBounds = true

	// 	  %7 = load i8*, i8** %6
	activateTaskVal7 := activateTaskIfThenIBlock.NewLoad(activateTaskVal6)

	// 	  %8 = bitcast i8* %7 to void (i8*)*
	activateTaskVal8 := activateTaskIfThenIBlock.NewBitCast(activateTaskVal7, types.NewPointer(types.NewFunc(types.Void, types.I8Ptr)))

	// 	  tail call fastcc void %8(i8* nonnull %0), !dbg !33
	activateTaskParam0v2 := *activateTaskParam0 // We use a copy of the parameter, as we don't want the "nonnull" showing up in the function declaration
	activateTaskParam0v2.Attrs = []ir.ParamAttribute{enum.ParamAttrNonNull}
	activateTaskCall2 := activateTaskIfThenIBlock.NewCall(activateTaskVal8, &activateTaskParam0v2)
	activateTaskCall2.Tail = enum.TailTail
	activateTaskCall2.CallingConv = enum.CallingConvFast
	activateTaskCall2.Metadata = append(activateTaskCall2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["33"]})

	// 	  br label %if.then, !dbg !34
	activateTaskBr3 := activateTaskIfThenIBlock.NewBr(activateTaskIfThenBlock)
	activateTaskBr3.Metadata = append(activateTaskBr3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["34"]})

	// 	if.done3.i:                                       ; preds = %if.done
	// 	%9 = load i8*, i8** @runtime.runqueueBack, align 4, !dbg !35
	activateTaskVal9 := activateTaskIfDone3IBlock.NewLoad(globalRunQueueBack)
	activateTaskVal9.Align = 4
	activateTaskVal9.Metadata = append(activateTaskVal9.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["35"]})

	// 	%10 = icmp eq i8* %9, null, !dbg !36
	activateTaskVal10 := activateTaskIfDone3IBlock.NewICmp(enum.IPredEQ, activateTaskVal9, null)
	activateTaskVal10.Metadata = append(activateTaskVal10.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["36"]})

	// 	br i1 %10, label %if.then4.i, label %store.next.i, !dbg !32
	activateTaskBr4 := activateTaskIfDone3IBlock.NewCondBr(activateTaskVal10, activateTaskIfThen4IBlock, activateTaskStoreNextIBlock)
	activateTaskBr4.Metadata = append(activateTaskBr4.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["32"]})

	// 	if.then4.i:                                       ; preds = %if.done3.i
	// 	store i8* %0, i8** @runtime.runqueueBack, align 4, !dbg !37
	activateTaskStore0 := activateTaskIfThen4IBlock.NewStore(activateTaskParam0, globalRunQueueBack)
	activateTaskStore0.Align = 4
	activateTaskStore0.Metadata = append(activateTaskStore0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["37"]})

	// 	store i8* %0, i8** @runtime.runqueueFront, align 4, !dbg !38
	activateTaskStore1 := activateTaskIfThen4IBlock.NewStore(activateTaskParam0, globalRunQueueFront)
	activateTaskStore1.Align = 4
	activateTaskStore1.Metadata = append(activateTaskStore1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["38"]})

	// 	br label %if.then, !dbg !32
	activateTaskBr5 := activateTaskIfThen4IBlock.NewBr(activateTaskIfThenBlock)
	activateTaskBr5.Metadata = append(activateTaskBr5.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["32"]})

	// 	store.next.i:                                     ; preds = %if.done3.i
	// 	call void @llvm.dbg.value(metadata i8* %9, metadata !39, metadata !DIExpression()), !dbg !42
	activateTaskCall3 := activateTaskStoreNextIBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: activateTaskVal9}, &metadata.Value{Value: dbgLoc["LocalVarT39"]}, &metadata.Value{Value: emptyExpr})
	activateTaskCall3.Metadata = append(activateTaskCall3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["42"]})

	// 	%11 = getelementptr inbounds i8, i8* %9, i32 8, !dbg !44
	activateTaskVal11 := activateTaskStoreNextIBlock.NewGetElementPtr(activateTaskVal9, constant.NewInt(types.I32, 8))
	activateTaskVal11.InBounds = true
	activateTaskVal11.Metadata = append(activateTaskVal11.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["44"]})

	// 	%12 = bitcast i8* %11 to i8**, !dbg !45
	activateTaskVal12 := activateTaskStoreNextIBlock.NewBitCast(activateTaskVal11, types.NewPointer(types.I8Ptr))
	activateTaskVal12.Metadata = append(activateTaskVal12.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["45"]})

	// 	store i8* %0, i8** %12, align 4, !dbg !45
	activateTaskStore2 := activateTaskStoreNextIBlock.NewStore(activateTaskParam0, activateTaskVal12)
	activateTaskStore2.Align = 4
	activateTaskStore2.Metadata = append(activateTaskStore2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["45"]})

	// 	store i8* %0, i8** @runtime.runqueueBack, align 4, !dbg !46
	activateTaskStore3 := activateTaskStoreNextIBlock.NewStore(activateTaskParam0, globalRunQueueBack)
	activateTaskStore3.Align = 4
	activateTaskStore3.Metadata = append(activateTaskStore3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["46"]})

	// 	br label %if.then, !dbg !32
	activateTaskBr6 := activateTaskStoreNextIBlock.NewBr(activateTaskIfThenBlock)
	activateTaskBr6.Metadata = append(activateTaskBr6.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["32"]})
	// 	}

	// 	; Function Attrs: optsize
	// 	define internal fastcc void @runtime.putchar(i8) unnamed_addr #0 section ".text.runtime.putchar" !dbg !111 {
	putCharFunc.Linkage = enum.LinkageInternal
	putCharFunc.CallingConv = enum.CallingConvFast
	putCharFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	putCharFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	putCharFunc.Section = ".text.runtime.putchar"
	putCharFunc.Metadata = append(putCharFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.putchar"]})

	// 	entry:
	putCharEntryBlock := putCharFunc.NewBlock("entry")

	// 	%stackalloc.alloca = alloca [1 x i32], align 4, !dbg !116
	putCharAlloc1 := putCharEntryBlock.NewAlloca(types.NewArray(1, types.I32))
	putCharAlloc1.LocalName = "stackalloc.alloca"
	putCharAlloc1.Align = 4
	putCharAlloc1.Metadata = append(putCharAlloc1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["116"]})

	// 	%.fca.0.gep = getelementptr inbounds [1 x i32], [1 x i32]* %stackalloc.alloca, i32 0, i32 0, !dbg !116
	putCharVal0 := putCharEntryBlock.NewGetElementPtr(putCharAlloc1, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, 0))
	putCharVal0.InBounds = true
	putCharVal0.LocalName = ".fca.0.gep"
	putCharVal0.Metadata = append(putCharVal0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["116"]})

	// 	store i32 0, i32* %.fca.0.gep, align 4, !dbg !116
	putCharStor0 := putCharEntryBlock.NewStore(zero, putCharVal0)
	putCharStor0.Align = 4
	putCharStor0.Metadata = append(putCharStor0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["116"]})

	// 	%stackalloc = bitcast [1 x i32]* %stackalloc.alloca to i8*, !dbg !116
	putCharBitCast0 := putCharEntryBlock.NewBitCast(putCharAlloc1, types.I8Ptr)
	putCharBitCast0.LocalName = "stackalloc"
	putCharBitCast0.Metadata = append(putCharBitCast0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["116"]})

	// 	call void @llvm.dbg.value(metadata i8 %0, metadata !115, metadata !DIExpression()), !dbg !116
	putCharCall0 := putCharEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: putCharParam0}, &metadata.Value{Value: dbgLoc["c115"]}, &metadata.Value{Value: emptyExpr})
	putCharCall0.Metadata = append(putCharCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["116"]})

	// 	store i8 %0, i8* %stackalloc, align 4, !dbg !117
	putCharStor1 := putCharEntryBlock.NewStore(putCharParam0, putCharBitCast0)
	putCharStor1.Align = 4
	putCharStor1.Metadata = append(putCharStor1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["117"]})

	// 	%1 = load i32, i32* @runtime.stdout, align 4, !dbg !118
	putCharLoad0 := putCharEntryBlock.NewLoad(globalStdout)
	putCharLoad0.Align = 4
	putCharLoad0.Metadata = append(putCharLoad0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["118"]})

	// 	%2 = call i32 @resource_write(i32 %1, i8* nonnull %stackalloc, i32 1), !dbg !119
	// TODO: Get the "nonnull" to emit
	//       https://github.com/llir/llvm/issues/88
	putCharCall1 := putCharEntryBlock.NewCall(resWriFunc, putCharLoad0, putCharBitCast0, one)
	putCharCall1.Metadata = append(putCharCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["119"]})

	// 	ret void, !dbg !117
	putCharRet := putCharEntryBlock.NewRet(nil)
	putCharRet.Metadata = append(putCharRet.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["117"]})
	// 	}

	// 	; Function Attrs: optsize
	// 	define internal fastcc void @runtime.printstring(i8* nocapture readonly, i32) unnamed_addr #0 section ".text.runtime.printstring" !dbg !95 {
	prtStrParam0.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture, enum.ParamAttrReadOnly}
	prtStrFunc.Linkage = enum.LinkageInternal
	prtStrFunc.CallingConv = enum.CallingConvFast
	prtStrFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	prtStrFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	prtStrFunc.Section = ".text.runtime.printstring"
	prtStrFunc.Metadata = append(prtStrFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.printstring"]})

	// Create the blocks
	prtStrEntryBlock := prtStrFunc.NewBlock("entry")
	prtStrForLoopBlock := prtStrFunc.NewBlock("for.loop")
	prtStrForBodyBlock := prtStrFunc.NewBlock("for.body")
	prtStrForDoneBlock := prtStrFunc.NewBlock("for.done")

	// 	entry:
	// 	call void @llvm.dbg.value(metadata i8* %0, metadata !103, metadata !DIExpression(DW_OP_LLVM_fragment, 0, 32)), !dbg !104
	prtStrCall0DI := &metadata.DIExpression{
		MetadataID: -1,
		Fields:     []metadata.DIExpressionField{enum.DwarfOpLLVMFragment, metadata.UintLit(0), metadata.UintLit(32)},
	}
	prtStrCall0 := prtStrEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: prtStrParam0}, &metadata.Value{Value: dbgLoc["s"]}, &metadata.Value{Value: prtStrCall0DI})
	prtStrCall0.Metadata = append(prtStrCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["104"]})

	// 	call void @llvm.dbg.value(metadata i32 %1, metadata !103, metadata !DIExpression(DW_OP_LLVM_fragment, 32, 32)), !dbg !104
	prtStrCall1DI := &metadata.DIExpression{
		MetadataID: -1,
		Fields:     []metadata.DIExpressionField{enum.DwarfOpLLVMFragment, metadata.UintLit(32), metadata.UintLit(32)},
	}
	prtStrCall1 := prtStrEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: prtStrParam1}, &metadata.Value{Value: dbgLoc["s"]}, &metadata.Value{Value: prtStrCall1DI})
	prtStrCall1.Metadata = append(prtStrCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["104"]})

	// 	br label %for.loop, !dbg !105
	prtStrBr0 := prtStrEntryBlock.NewBr(prtStrForLoopBlock)
	prtStrBr0.Metadata = append(prtStrBr0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["105"]})

	// 	for.loop:                                         ; preds = %for.body, %entry
	// 	%2 = phi i32 [ 0, %entry ], [ %6, %for.body ], !dbg !106
	prtStrVal6 := &ir.InstAdd{}
	prtStrVal6.LocalID = 6 // Manually set this here, as automatic calculation doesn't work when done this way
	prtStrVal2 := prtStrForLoopBlock.NewPhi(ir.NewIncoming(zero, prtStrEntryBlock), ir.NewIncoming(prtStrVal6, prtStrForBodyBlock))
	prtStrVal2.Metadata = append(prtStrVal2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["106"]})

	// 	%3 = icmp slt i32 %2, %1, !dbg !107
	prtStrVal3 := prtStrForLoopBlock.NewICmp(enum.IPredSLT, prtStrVal2, prtStrParam1)
	prtStrVal3.Metadata = append(prtStrVal3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["107"]})

	// 	br i1 %3, label %for.body, label %for.done, !dbg !105
	prtStrBr1 := prtStrForLoopBlock.NewCondBr(prtStrVal3, prtStrForBodyBlock, prtStrForDoneBlock)
	prtStrBr1.Metadata = append(prtStrBr1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["105"]})

	// 	for.body:                                         ; preds = %for.loop
	// 	%4 = getelementptr inbounds i8, i8* %0, i32 %2, !dbg !108
	prtStrVal4 := prtStrForBodyBlock.NewGetElementPtr(prtStrParam0, prtStrVal2)
	prtStrVal4.InBounds = true
	prtStrVal4.Metadata = append(prtStrVal4.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["108"]})

	// 	%5 = load i8, i8* %4, align 1, !dbg !108
	prtStrVal5 := prtStrForBodyBlock.NewLoad(prtStrVal4)
	prtStrVal5.Align = 1
	prtStrVal5.Metadata = append(prtStrVal5.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["108"]})

	// 	tail call fastcc void @runtime.putchar(i8 %5), !dbg !109
	prtStrCall2 := prtStrForBodyBlock.NewCall(putCharFunc, prtStrVal5)
	prtStrCall2.Tail = enum.TailTail
	prtStrCall2.CallingConv = enum.CallingConvFast
	prtStrCall2.Metadata = append(prtStrCall2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["109"]})

	// 	%6 = add nuw i32 %2, 1, !dbg !110
	prtStrVal6 = prtStrForBodyBlock.NewAdd(prtStrVal2, constant.NewInt(types.I32, 1))
	prtStrVal6.OverflowFlags = []enum.OverflowFlag{enum.OverflowFlagNUW}
	prtStrVal6.Metadata = append(prtStrVal6.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["110"]})

	// 	br label %for.loop, !dbg !105
	prtStrBr2 := prtStrForBodyBlock.NewBr(prtStrForLoopBlock)
	prtStrBr2.Metadata = append(prtStrBr2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["105"]})

	// 	for.done:                                         ; preds = %for.loop
	// 	ret void, !dbg !105
	prtStrRet := prtStrForDoneBlock.NewRet(nil)
	prtStrRet.Metadata = append(prtStrRet.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["105"]})
	// 	}

	// 	; Function Attrs: optsize
	// 	define internal fastcc void @runtime.printnl() unnamed_addr #0 section ".text.runtime.printnl" !dbg !90 {
	prtNlFunc.Linkage = enum.LinkageInternal
	prtNlFunc.CallingConv = enum.CallingConvFast
	prtNlFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	prtNlFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	prtNlFunc.Section = ".text.runtime.printnl"
	prtNlFunc.Metadata = append(prtNlFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.printnl"]})

	// 	entry:
	prtNlEntryBlock := prtNlFunc.NewBlock("entry")

	// 	tail call fastcc void @runtime.putchar(i8 13), !dbg !92
	prtNlCall1 := prtNlEntryBlock.NewCall(putCharFunc, thirteen)
	prtNlCall1.Tail = enum.TailTail
	prtNlCall1.CallingConv = enum.CallingConvFast
	prtNlCall1.Metadata = append(prtNlCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["92"]})

	// 	tail call fastcc void @runtime.putchar(i8 10), !dbg !93
	prtNlCall2 := prtNlEntryBlock.NewCall(putCharFunc, ten)
	prtNlCall2.Tail = enum.TailTail
	prtNlCall2.CallingConv = enum.CallingConvFast
	prtNlCall2.Metadata = append(prtNlCall2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["93"]})

	// 	ret void, !dbg !94
	prtNlRet := prtNlEntryBlock.NewRet(nil)
	prtNlRet.Metadata = append(prtNlRet.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["94"]})
	// 	}

	// 	; Function Attrs: optsize
	// 	define void @cwa_main() local_unnamed_addr #0 section ".text.cwa_main" !dbg !47 {
	cwaMainFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	cwaMainFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	cwaMainFunc.Section = ".text.cwa_main"
	cwaMainFunc.Metadata = append(cwaMainFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.cwa_main"]})

	// 	entry:
	cwaMainEntryBlock := cwaMainFunc.NewBlock("entry")

	// 	%0 = tail call i32 @io_get_stdout(), !dbg !48
	cwaVal0 := cwaMainEntryBlock.NewCall(IOGetStdoutFunc)
	cwaVal0.Metadata = append(cwaVal0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["48"]})

	// 	store i32 %0, i32* @runtime.stdout, align 4, !dbg !48
	cwaStore0 := cwaMainEntryBlock.NewStore(cwaVal0, globalStdout)
	cwaStore0.Align = 4
	cwaStore0.Metadata = append(cwaStore0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["48"]})

	// 	tail call fastcc void @runtime.printstring(i8* getelementptr inbounds ([12 x i8], [12 x i8]* @"main.go.main$string", i32 0, i32 0), i32 12), !dbg !50
	// TODO: See if I can figure out how to get the resulting IR to be on one line, instead of split over two
	cwaCall0 := cwaMainEntryBlock.NewCall(prtStrFunc, cwaMainEntryBlock.NewGetElementPtr(globalMain, zero, zero), twelve)
	cwaCall0.Tail = enum.TailTail
	cwaCall0.CallingConv = enum.CallingConvFast
	cwaCall0.Metadata = append(cwaCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["50"]})

	// 	tail call fastcc void @runtime.printnl(), !dbg !50
	cwaCall1 := cwaMainEntryBlock.NewCall(prtNlFunc)
	cwaCall1.Tail = enum.TailTail
	cwaCall1.CallingConv = enum.CallingConvFast
	cwaCall1.Metadata = append(cwaCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["50"]})

	// 	ret void, !dbg !53
	cwaRet := cwaMainEntryBlock.NewRet(nil)
	cwaRet.Metadata = append(cwaRet.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["53"]})
	// 	}

	// 	; Function Attrs: noreturn optsize
	// 	define internal fastcc void @runtime.runtimePanic() unnamed_addr #1 section ".text.runtime.runtimePanic" !dbg !122 {
	runtimePanicFunc.Linkage = enum.LinkageInternal
	runtimePanicFunc.CallingConv = enum.CallingConvFast
	runtimePanicFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	runtimePanicFunc.FuncAttrs = []ir.FuncAttribute{&atGrp1}
	runtimePanicFunc.Section = ".text.runtime.runtimePanic"
	runtimePanicFunc.Metadata = append(runtimePanicFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.runtimePanic"]})

	// 	entry:
	runtimePanicEntryBlock := runtimePanicFunc.NewBlock("entry")

	// 	call void @llvm.dbg.value(metadata i8* getelementptr inbounds ([23 x i8], [23 x i8]* @"runtime.nilPanic$string", i32 0, i32 0), metadata !124, metadata !DIExpression(DW_OP_LLVM_fragment, 0, 32)), !dbg !125
	// TODO: See if I can figure out how to get the resulting IR to be on one line, instead of split over two (there are several locations like this)
	// TODO: Also check if the "inbounds" is being emitted when done in this way
	runtimePanicCall0 := runtimePanicEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: runtimePanicEntryBlock.NewGetElementPtr(globalNilPanic, zero, zero)}, &metadata.Value{Value: dbgLoc["msg"]}, &metadata.Value{Value: prtStrCall0DI})
	runtimePanicCall0.Metadata = append(runtimePanicCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["125"]})

	// 	call void @llvm.dbg.value(metadata i32 23, metadata !124, metadata !DIExpression(DW_OP_LLVM_fragment, 32, 32)), !dbg !125
	runtimePanicCall1 := runtimePanicEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: constant.NewInt(types.I32, 23)}, &metadata.Value{Value: dbgLoc["msg"]}, &metadata.Value{Value: prtStrCall1DI})
	runtimePanicCall1.Metadata = append(runtimePanicCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["125"]})

	// 	tail call fastcc void @runtime.printstring(i8* getelementptr inbounds ([22 x i8], [22 x i8]* @"runtime.runtimePanic$string", i32 0, i32 0), i32 22), !dbg !126
	runtimePanicCall2 := runtimePanicEntryBlock.NewCall(prtStrFunc, runtimePanicEntryBlock.NewGetElementPtr(globalNilPanic, zero, zero), constant.NewInt(types.I32, 22))
	runtimePanicCall2.Tail = enum.TailTail
	runtimePanicCall2.CallingConv = enum.CallingConvFast
	runtimePanicCall2.Metadata = append(runtimePanicCall2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["126"]})

	// 	tail call fastcc void @runtime.printstring(i8* getelementptr inbounds ([23 x i8], [23 x i8]* @"runtime.nilPanic$string", i32 0, i32 0), i32 23), !dbg !127
	runtimePanicCall3 := runtimePanicEntryBlock.NewCall(prtStrFunc, runtimePanicEntryBlock.NewGetElementPtr(globalNilPanic, zero, zero), constant.NewInt(types.I32, 23))
	runtimePanicCall3.Tail = enum.TailTail
	runtimePanicCall3.CallingConv = enum.CallingConvFast
	runtimePanicCall3.Metadata = append(runtimePanicCall3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["127"]})

	// 	tail call fastcc void @runtime.printnl(), !dbg !127
	runtimePanicCall4 := runtimePanicEntryBlock.NewCall(prtNlFunc)
	runtimePanicCall4.Tail = enum.TailTail
	runtimePanicCall4.CallingConv = enum.CallingConvFast
	runtimePanicCall4.Metadata = append(runtimePanicCall4.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["127"]})

	// 	tail call void @llvm.trap() #5, !dbg !128
	runtimePanicCall5 := runtimePanicEntryBlock.NewCall(LLVMTrapFunc)
	runtimePanicCall5.Tail = enum.TailTail
	runtimePanicCall5.FuncAttrs = []ir.FuncAttribute{&atGrp5}
	runtimePanicCall5.Metadata = append(runtimePanicCall5.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["128"]})

	// 	unreachable, !dbg !131
	runtimePanicUnreach0 := runtimePanicEntryBlock.NewUnreachable()
	runtimePanicUnreach0.Metadata = append(runtimePanicUnreach0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["131"]})
	// 	}

	// 	; Function Attrs: noreturn optsize
	// 	define internal fastcc void @runtime.nilPanic() unnamed_addr #1 section ".text.runtime.nilPanic" !dbg !87 {
	nilPanicFunc.Linkage = enum.LinkageInternal
	nilPanicFunc.CallingConv = enum.CallingConvFast
	nilPanicFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	nilPanicFunc.FuncAttrs = []ir.FuncAttribute{&atGrp1}
	nilPanicFunc.Section = ".text.runtime.nilPanic"
	nilPanicFunc.Metadata = append(nilPanicFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.nilPanic"]})

	// 	entry:
	nilPanicEntryBlock := nilPanicFunc.NewBlock("entry")

	// 	tail call fastcc void @runtime.runtimePanic(), !dbg !89
	nilPanicCall0 := nilPanicEntryBlock.NewCall(runtimePanicFunc)
	nilPanicCall0.Tail = enum.TailTail
	nilPanicCall0.CallingConv = enum.CallingConvFast
	nilPanicCall0.Metadata = append(nilPanicCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["89"]})

	// 	unreachable
	nilPanicEntryBlock.NewUnreachable()
	// 	}

	// 	; Function Attrs: noreturn optsize
	// 	define internal fastcc void @runtime.getFuncPtr() unnamed_addr #1 section ".text.runtime.getFuncPtr" !dbg !54 {
	getFuncPtrFunc.Linkage = enum.LinkageInternal
	getFuncPtrFunc.CallingConv = enum.CallingConvFast
	getFuncPtrFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	getFuncPtrFunc.FuncAttrs = []ir.FuncAttribute{&atGrp1}
	getFuncPtrFunc.Section = ".text.runtime.getFuncPtr"
	getFuncPtrFunc.Metadata = append(getFuncPtrFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.getFuncPtr"]})

	// 	entry:
	getFuncPtrEntryBlock := getFuncPtrFunc.NewBlock("entry")

	// 	call void @llvm.dbg.value(metadata i8* null, metadata !67, metadata !DIExpression(DW_OP_LLVM_fragment, 0, 32)), !dbg !69
	getFuncPtrCall0 := getFuncPtrEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: &i8nullptr}, &metadata.Value{Value: dbgLoc["LocalVarVal"]}, &metadata.Value{Value: prtStrCall0DI})
	getFuncPtrCall0.Metadata = append(getFuncPtrCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["69"]})

	// 	call void @llvm.dbg.value(metadata i32 0, metadata !67, metadata !DIExpression(DW_OP_LLVM_fragment, 32, 32)), !dbg !69
	getFuncPtrCall1 := getFuncPtrEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: zero}, &metadata.Value{Value: dbgLoc["LocalVarVal"]}, &metadata.Value{Value: prtStrCall1DI})
	getFuncPtrCall1.Metadata = append(getFuncPtrCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["69"]})

	// 	call void @llvm.dbg.value(metadata i8* undef, metadata !68, metadata !DIExpression()), !dbg !69
	getFuncPtrCall2 := getFuncPtrEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: &i8nullptr}, &metadata.Value{Value: dbgLoc["LocalVarSignature"]}, &metadata.Value{Value: emptyExpr})
	getFuncPtrCall2.Metadata = append(getFuncPtrCall2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["69"]})

	// 	tail call fastcc void @runtime.nilPanic(), !dbg !70
	getFuncPtrCall3 := getFuncPtrEntryBlock.NewCall(nilPanicFunc)
	getFuncPtrCall3.Tail = enum.TailTail
	getFuncPtrCall3.CallingConv = enum.CallingConvFast
	getFuncPtrCall3.Metadata = append(getFuncPtrCall3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["70"]})

	// 	unreachable, !dbg !70
	getFuncPtrUnreach0 := getFuncPtrEntryBlock.NewUnreachable()
	getFuncPtrUnreach0.Metadata = append(getFuncPtrUnreach0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["70"]})
	// 	}

	// 	; Function Attrs: optsize
	// 	define i8* @memset(i8* nocapture returned, i8, i32) local_unnamed_addr #0 section ".text.memset" !dbg !71 {
	memSetParam0.Attrs = []ir.ParamAttribute{enum.ParamAttrNoCapture, enum.ParamAttrReturned}
	memSetFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	memSetFunc.FuncAttrs = []ir.FuncAttribute{&atGrp0}
	memSetFunc.Section = ".text.memset"
	memSetFunc.Metadata = append(memSetFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.memset"]})

	// 	entry:
	memSetEntryBlock := memSetFunc.NewBlock("entry")
	memSetForLoopBlock := memSetFunc.NewBlock("for.loop")
	memSetForBodyBlock := memSetFunc.NewBlock("for.body")
	memSetForDoneBlock := memSetFunc.NewBlock("for.done")
	memSetStoreNilBlock := memSetFunc.NewBlock("store.nil")
	memSetStoreNextBlock := memSetFunc.NewBlock("store.next")

	// 	call void @llvm.dbg.value(metadata i8* %0, metadata !76, metadata !DIExpression()), !dbg !79
	memSetCall0 := memSetEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: memSetParam0}, &metadata.Value{Value: dbgLoc["ptr"]}, &metadata.Value{Value: emptyExpr})
	memSetCall0.Metadata = append(memSetCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["79"]})

	// 	call void @llvm.dbg.value(metadata i8 %1, metadata !77, metadata !DIExpression()), !dbg !79
	memSetFuncCall1 := memSetEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: memSetParam1}, &metadata.Value{Value: dbgLoc["c"]}, &metadata.Value{Value: emptyExpr})
	memSetFuncCall1.Metadata = append(memSetFuncCall1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["79"]})

	// 	call void @llvm.dbg.value(metadata i32 %2, metadata !78, metadata !DIExpression()), !dbg !79
	memSetFuncCall2 := memSetEntryBlock.NewCall(LLVMDbgValueFunc, &metadata.Value{Value: memSetParam2}, &metadata.Value{Value: dbgLoc["size"]}, &metadata.Value{Value: emptyExpr})
	memSetFuncCall2.Metadata = append(memSetFuncCall2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["79"]})

	// 	br label %for.loop, !dbg !80
	memSetFuncBr0 := memSetEntryBlock.NewBr(memSetForLoopBlock)
	memSetFuncBr0.Metadata = append(memSetFuncBr0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["80"]})

	// 	for.loop:                                         ; preds = %store.next, %entry
	// 	%3 = phi i32 [ 0, %entry ], [ %7, %store.next ], !dbg !81
	memSetFuncVal7 := &ir.InstAdd{}
	memSetFuncVal7.LocalID = 7 // Manually set this here, as automatic calculation doesn't work when done this way
	memSetFuncVal3 := memSetForLoopBlock.NewPhi(ir.NewIncoming(zero, memSetEntryBlock), ir.NewIncoming(memSetFuncVal7, memSetStoreNextBlock))
	memSetFuncVal3.Metadata = append(memSetFuncVal3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["81"]})

	// 	%4 = icmp ult i32 %3, %2, !dbg !82
	memSetFuncVal4 := memSetForLoopBlock.NewICmp(enum.IPredULT, memSetFuncVal3, memSetParam2)
	memSetFuncVal4.Metadata = append(memSetFuncVal4.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["82"]})

	// 	br i1 %4, label %for.body, label %for.done, !dbg !80
	memSetFuncBr1 := memSetForLoopBlock.NewCondBr(memSetFuncVal4, memSetForBodyBlock, memSetForDoneBlock)
	memSetFuncBr1.Metadata = append(memSetFuncBr1.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["80"]})

	// 	for.body:                                         ; preds = %for.loop
	// 	%5 = getelementptr inbounds i8, i8* %0, i32 %3, !dbg !83
	memSetFuncVal5 := memSetForBodyBlock.NewGetElementPtr(memSetParam0, memSetFuncVal3)
	memSetFuncVal5.InBounds = true
	memSetFuncVal5.Metadata = append(memSetFuncVal5.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["83"]})

	// 	%6 = icmp eq i8* %5, null, !dbg !84
	memSetFuncVal6 := memSetForBodyBlock.NewICmp(enum.IPredEQ, memSetFuncVal5, null)
	memSetFuncVal6.Metadata = append(memSetFuncVal6.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["84"]})

	// 	br i1 %6, label %store.nil, label %store.next, !dbg !84
	memSetFuncBr2 := memSetForBodyBlock.NewCondBr(memSetFuncVal6, memSetStoreNilBlock, memSetStoreNextBlock)
	memSetFuncBr2.Metadata = append(memSetFuncBr2.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["84"]})

	// 	for.done:                                         ; preds = %for.loop
	// 	ret i8* %0, !dbg !85
	memSetRet0 := memSetForDoneBlock.NewRet(memSetParam0)
	memSetRet0.Metadata = append(memSetRet0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["85"]})

	// 	store.nil:                                        ; preds = %for.body
	// 	tail call fastcc void @runtime.nilPanic(), !dbg !84
	memSetCall3 := memSetStoreNilBlock.NewCall(nilPanicFunc)
	memSetCall3.Tail = enum.TailTail
	memSetCall3.CallingConv = enum.CallingConvFast
	memSetCall3.Metadata = append(memSetCall3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["84"]})

	// 	unreachable, !dbg !84
	memSetUnreach0 := memSetStoreNilBlock.NewUnreachable()
	memSetUnreach0.Metadata = append(memSetUnreach0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["84"]})

	// 	store.next:                                       ; preds = %for.body
	// 	store i8 %1, i8* %5, align 1, !dbg !84
	memSetStore0 := memSetStoreNextBlock.NewStore(memSetParam1, memSetFuncVal5)
	memSetStore0.Align = 1
	memSetStore0.Metadata = append(memSetStore0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["84"]})

	// 	%7 = add i32 %3, 1, !dbg !86
	memSetFuncVal7 = memSetStoreNextBlock.NewAdd(memSetFuncVal3, one)
	memSetFuncVal7.Metadata = append(memSetFuncVal7.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["86"]})

	// 	br label %for.loop, !dbg !80
	memSetFuncBr3 := memSetStoreNextBlock.NewBr(memSetForLoopBlock)
	memSetFuncBr3.Metadata = append(memSetFuncBr3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["80"]})
	// 	}

	// 	; Function Attrs: noreturn optsize
	// 	define void @resume() local_unnamed_addr #1 section ".text.resume" !dbg !120 {
	resumeFunc.UnnamedAddr = enum.UnnamedAddrUnnamedAddr
	resumeFunc.FuncAttrs = []ir.FuncAttribute{&atGrp1}
	resumeFunc.Section = ".text.resume"
	resumeFunc.Metadata = append(resumeFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["runtime.resume"]})

	// 	entry:
	resumeFuncEntryBlock := resumeFunc.NewBlock("entry")

	// 	tail call fastcc void @runtime.getFuncPtr(), !dbg !121
	resumeFuncCall0 := resumeFuncEntryBlock.NewCall(getFuncPtrFunc)
	resumeFuncCall0.Tail = enum.TailTail
	resumeFuncCall0.CallingConv = enum.CallingConvFast
	resumeFuncCall0.Metadata = append(resumeFuncCall0.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["121"]})

	// 	unreachable
	resumeFuncEntryBlock.NewUnreachable()
	// 	}

	// Print LLVM IR assembly to standard output.
	// m.String()
	fmt.Println(m)
}

func addMetadata(m *ir.Module) {

	// Convenience constants.
	var (
		one   = constant.NewInt(types.I32, 1)
		three = constant.NewInt(types.I32, 3)
		four  = constant.NewInt(types.I32, 4)
	)

	// * Unnamed metadata definitions *

	// 	!0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "TinyGo", isOptimized: true, runtimeVersion: 0, emissionKind: FullDebug, enums: !2)
	diCompileUnit := &metadata.DICompileUnit{
		MetadataID:     -1,
		Distinct:       true,
		Language:       enum.DwarfLangC99,
		Producer:       "TinyGo",
		IsOptimized:    true,
		RuntimeVersion: 0,
		EmissionKind:   enum.EmissionKindFullDebug,
	}

	// 	!1 = !DIFile(filename: "main.go", directory: "")
	diFileMainGo1 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "main.go",
		Directory:  "",
	}
	diCompileUnit.File = diFileMainGo1

	// 	!2 = !{}
	emptyTuple2 := &metadata.Tuple{
		MetadataID: -1,
	}
	diCompileUnit.Enums = emptyTuple2

	// 	!3 = !{i32 1, !"Debug Info Version", i32 3}
	debugInfoVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{one, &metadata.String{Value: "Debug Info Version"}, three},
	}

	// 	!4 = !{i32 1, !"Dwarf Version", i32 4}
	dwarfVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{one, &metadata.String{Value: "Dwarf Version"}, four},
	}

	// 	!5 = distinct !DISubprogram(name: "runtime._start", linkageName: "_start", scope: !6, file: !6, line: 26, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogram_start := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime._start",
		LinkageName:   "_start",
		Line:          26,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["_start"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogram_start))

	// 	!6 = !DIFile(filename: "runtime_wasm.go", directory: "../../../runtime")
	diFileRuntimeWasmGo6 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "runtime_wasm.go",
		Directory:  "../../../runtime",
	}
	diSubprogram_start.Scope = diFileRuntimeWasmGo6
	diSubprogram_start.File = diFileRuntimeWasmGo6

	// 	!7 = !DISubroutineType(types: !2)
	diSubroutineType7 := &metadata.DISubroutineType{
		MetadataID: -1,
		Types:      emptyTuple2,
	}
	diSubprogram_start.Type = diSubroutineType7

	// 	!8 = !DILocation(line: 11, column: 6, scope: !9, inlinedAt: !11)
	dbgLoc["8"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       11,
		Column:     6,
	}

	// 	!9 = distinct !DISubprogram(name: "runtime.initAll", linkageName: "runtime.initAll", scope: !10, file: !10, line: 11, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramRuntimeInitAll9 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime._start",
		LinkageName:   "runtime.initAll",
		Line:          11,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
		Type:          diSubroutineType7,
	}
	dbgLoc["runtime.initAll"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeInitAll9))
	dbgLoc["8"].Scope = diSubprogramRuntimeInitAll9

	// 	!10 = !DIFile(filename: "runtime.go", directory: "../../../runtime")
	diFileRuntimeGo10 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "runtime.go",
		Directory:  "../../../runtime",
	}
	diSubprogramRuntimeInitAll9.Scope = diFileRuntimeGo10
	diSubprogramRuntimeInitAll9.File = diFileRuntimeGo10

	// 	!11 = distinct !DILocation(line: 27, column: 9, scope: !5)
	dbgLoc["11"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       27,
		Column:     9,
		Scope:      diSubprogram_start,
	}
	dbgLoc["8"].InlinedAt = dbgLoc["11"]

	// 	!12 = !DILocation(line: 0, scope: !5)
	dbgLoc["12"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogram_start,
	}

	// 	!13 = distinct !DISubprogram(name: "runtime.activateTask", linkageName: "runtime.activateTask", scope: !14, file: !14, line: 106, type: !15, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !20)
	diSubprogramRuntimeActivateTask13 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.activateTask",
		LinkageName: "runtime.activateTask",
		Line:        106,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.activateTask"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeActivateTask13))

	// 	!14 = !DIFile(filename: ".go", directory: "../../../runtime")
	diFileSchedulerGo14 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "scheduler.go",
		Directory:  "../../../runtime",
	}
	diSubprogramRuntimeActivateTask13.Scope = diFileSchedulerGo14
	diSubprogramRuntimeActivateTask13.File = diFileSchedulerGo14

	// 	!15 = !DISubroutineType(types: !16)
	diSubroutine15 := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramRuntimeActivateTask13.Type = diSubroutine15

	// 	!16 = !{!17}
	diTuple16 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubroutine15.Types = diTuple16

	// 	!17 = !DIDerivedType(tag: DW_TAG_pointer_type, baseType: !18, size: 32, align: 32, dwarfAddressSpace: 0)
	diDerived17 := &metadata.DIDerivedType{
		MetadataID:        -1,
		Tag:               enum.DwarfTagPointerType,
		Size:              32,
		Align:             32,
		DwarfAddressSpace: 0,
	}
	diTuple16.Fields = append(diTuple16.Fields, diDerived17)

	// 	!18 = !DIDerivedType(tag: DW_TAG_typedef, name: "runtime.coroutine", baseType: !19)
	diDerived18 := &metadata.DIDerivedType{
		MetadataID: -1,
		Tag:        enum.DwarfTagTypedef,
		Name:       "runtime.coroutine",
	}
	diDerived17.BaseType = diDerived18

	// 	!19 = !DIBasicType(name: "uint8", size: 8, encoding: DW_ATE_unsigned)
	diBasicType19 := &metadata.DIBasicType{
		MetadataID: -1,
		Name:       "uint8",
		Size:       8,
		Encoding:   enum.DwarfAttEncodingUnsigned,
	}
	diDerived18.BaseType = diBasicType19

	// 	!20 = !{!21}
	diTuple20 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubprogramRuntimeActivateTask13.RetainedNodes = diTuple20

	// 	!21 = !DILocalVariable(name: "task", arg: 1, scope: !13, file: !14, line: 106, type: !17)
	diLocalVarTask21 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "task",
		Arg:        1,
		Scope:      diSubprogramRuntimeActivateTask13,
		File:       diFileSchedulerGo14,
		Line:       106,
		Type:       diDerived17,
	}
	dbgLoc["LocalVarTask"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarTask21))
	diTuple20.Fields = append(diTuple20.Fields, diLocalVarTask21)

	// 	!22 = !DILocation(line: 106, column: 6, scope: !13)
	dbgLoc["22"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       106,
		Column:     6,
		Scope:      diSubprogramRuntimeActivateTask13,
	}

	// 	!23 = !DILocation(line: 107, column: 10, scope: !13)
	dbgLoc["23"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       107,
		Column:     10,
		Scope:      diSubprogramRuntimeActivateTask13,
	}

	// 	!24 = !DILocation(line: 0, scope: !13)
	dbgLoc["24"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Column:     13,
		Scope:      diSubprogramRuntimeActivateTask13,
	}

	// 	!25 = !DILocation(line: 108, column: 3, scope: !13)
	dbgLoc["25"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       108,
		Column:     3,
		Scope:      diSubprogramRuntimeActivateTask13,
	}

	// 	!26 = !DILocalVariable(name: "t", arg: 1, scope: !27, file: !14, line: 137, type: !17)
	diLocalVarT26 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "t",
		Arg:        1,
		File:       diFileSchedulerGo14,
		Line:       137,
		Type:       diDerived17,
	}
	dbgLoc["LocalVarT"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarT26))

	// 	!27 = distinct !DISubprogram(name: "runtime.runqueuePushBack", linkageName: "runtime.runqueuePushBack", scope: !14, file: !14, line: 137, type: !15, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !28)
	diSubprogramRuntimeRunqueuePushBack27 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.runqueuePushBack",
		LinkageName: "runtime.runqueuePushBack",
		Scope:       diFileSchedulerGo14,
		File:        diFileSchedulerGo14,
		Line:        137,
		Type:        diSubroutine15,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.runqueuePushBack"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeRunqueuePushBack27))
	diLocalVarT26.Scope = diSubprogramRuntimeRunqueuePushBack27

	// 	!28 = !{!26}
	diTuple28 := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{diLocalVarT26},
	}
	diSubprogramRuntimeRunqueuePushBack27.RetainedNodes = diTuple28

	// 	!29 = !DILocation(line: 137, column: 6, scope: !27, inlinedAt: !30)
	dbgLoc["29"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       137,
		Column:     6,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
	}

	// 	!30 = distinct !DILocation(line: 111, column: 18, scope: !13)
	dbgLoc["30"] = &metadata.DILocation{
		MetadataID: -1,
		Distinct:   true,
		Line:       111,
		Column:     18,
		Scope:      diSubprogramRuntimeActivateTask13,
	}
	dbgLoc["29"].InlinedAt = dbgLoc["30"]

	// 	!31 = !DILocation(line: 138, column: 11, scope: !27, inlinedAt: !30)
	dbgLoc["31"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       138,
		Column:     11,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!32 = !DILocation(line: 0, scope: !27, inlinedAt: !30)
	dbgLoc["32"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!33 = !DILocation(line: 140, column: 12, scope: !27, inlinedAt: !30)
	dbgLoc["33"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       140,
		Column:     12,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!34 = !DILocation(line: 141, column: 3, scope: !27, inlinedAt: !30)
	dbgLoc["34"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       141,
		Column:     3,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!35 = !DILocation(line: 148, column: 5, scope: !27, inlinedAt: !30)
	dbgLoc["35"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       148,
		Column:     5,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!36 = !DILocation(line: 148, column: 18, scope: !27, inlinedAt: !30)
	dbgLoc["36"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       148,
		Column:     18,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!37 = !DILocation(line: 150, column: 3, scope: !27, inlinedAt: !30)
	dbgLoc["37"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       150,
		Column:     3,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!38 = !DILocation(line: 151, column: 3, scope: !27, inlinedAt: !30)
	dbgLoc["38"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       151,
		Column:     3,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!39 = !DILocalVariable(name: "t", arg: 1, scope: !40, file: !14, line: 46, type: !17)
	diLocalVarT39 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "t",
		Arg:        1,
		File:       diFileSchedulerGo14,
		Line:       46,
		Type:       diDerived17,
	}
	dbgLoc["LocalVarT39"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarT39))

	// 	!40 = distinct !DISubprogram(name: "(*runtime.coroutine).promise", linkageName: "(*runtime.coroutine).promise", scope: !14, file: !14, line: 46, type: !15, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !41)
	diSubprogramRuntimeCoroutinePromisePtr40 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "(*runtime.coroutine).promise",
		LinkageName: "(*runtime.coroutine).promise",
		Scope:       diFileSchedulerGo14,
		File:        diFileSchedulerGo14,
		Line:        46,
		Type:        diSubroutine15,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.runqueuePushBack"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeCoroutinePromisePtr40))
	diLocalVarT39.Scope = diSubprogramRuntimeCoroutinePromisePtr40

	// 	!41 = !{!39}
	diTuple41 := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{diLocalVarT39},
	}
	diSubprogramRuntimeCoroutinePromisePtr40.RetainedNodes = diTuple41

	// 	!42 = !DILocation(line: 46, column: 21, scope: !40, inlinedAt: !43)
	dbgLoc["42"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       46,
		Column:     21,
		Scope:      diSubprogramRuntimeCoroutinePromisePtr40,
	}

	// 	!43 = distinct !DILocation(line: 154, column: 42, scope: !27, inlinedAt: !30)
	dbgLoc["43"] = &metadata.DILocation{
		MetadataID: -1,
		Distinct:   true,
		Line:       154,
		Column:     42,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}
	dbgLoc["42"].InlinedAt = dbgLoc["43"]

	// 	!44 = !DILocation(line: 47, column: 32, scope: !40, inlinedAt: !43)
	dbgLoc["44"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       47,
		Column:     32,
		Scope:      diSubprogramRuntimeCoroutinePromisePtr40,
		InlinedAt:  dbgLoc["43"],
	}

	// 	!45 = !DILocation(line: 155, column: 19, scope: !27, inlinedAt: !30)
	dbgLoc["45"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       155,
		Column:     19,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!46 = !DILocation(line: 156, column: 3, scope: !27, inlinedAt: !30)
	dbgLoc["46"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       156,
		Column:     3,
		Scope:      diSubprogramRuntimeRunqueuePushBack27,
		InlinedAt:  dbgLoc["30"],
	}

	// 	!47 = distinct !DISubprogram(name: "runtime.cwa_main", linkageName: "cwa_main", scope: !6, file: !6, line: 31, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramRuntimeCwaMain47 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime.cwa_main",
		LinkageName:   "cwa_main",
		Scope:         diFileRuntimeWasmGo6,
		File:          diFileRuntimeWasmGo6,
		Line:          31,
		Type:          diSubroutineType7,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["runtime.cwa_main"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeCwaMain47))

	// 	!48 = !DILocation(line: 11, column: 6, scope: !9, inlinedAt: !49)
	dbgLoc["48"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       11,
		Column:     6,
		Scope:      diSubprogramRuntimeInitAll9,
	}

	// 	!49 = distinct !DILocation(line: 32, column: 9, scope: !47)
	dbgLoc["49"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       32,
		Column:     9,
		Scope:      diSubprogramRuntimeCwaMain47,
	}
	dbgLoc["48"].InlinedAt = dbgLoc["49"]

	// 	!50 = !DILocation(line: 4, column: 9, scope: !51, inlinedAt: !52)
	dbgLoc["50"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       4,
		Column:     9,
	}

	// 	!51 = distinct !DISubprogram(name: "main.go.main", linkageName: "main.go.main", scope: !1, file: !1, line: 3, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramMainGoMain51 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "main.go.main",
		LinkageName:   "main.go.main",
		Scope:         diFileMainGo1,
		File:          diFileMainGo1,
		Line:          3,
		Type:          diSubroutineType7,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["main.go.main"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramMainGoMain51))
	dbgLoc["50"].Scope = diSubprogramMainGoMain51

	// 	!52 = distinct !DILocation(line: 33, column: 10, scope: !47)
	dbgLoc["52"] = &metadata.DILocation{
		MetadataID: -1,
		Distinct:   true,
		Line:       33,
		Column:     10,
		Scope:      diSubprogramRuntimeCwaMain47,
	}
	dbgLoc["50"].InlinedAt = dbgLoc["52"]

	// 	!53 = !DILocation(line: 0, scope: !47)
	dbgLoc["53"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimeCwaMain47,
	}

	// 	!54 = distinct !DISubprogram(name: "runtime.getFuncPtr", linkageName: "runtime.getFuncPtr", scope: !55, file: !55, line: 26, type: !56, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !66)
	diSubprogramRuntimeGetFuncPtr54 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.getFuncPtr",
		LinkageName: "runtime.getFuncPtr",
		Line:        26,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.getFuncPtr"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeGetFuncPtr54))

	// 	!55 = !DIFile(filename: "func.go", directory: "../../../runtime")
	diFileFuncGo55 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "func.go",
		Directory:  "../../../runtime",
	}
	diSubprogramRuntimeGetFuncPtr54.Scope = diFileFuncGo55
	diSubprogramRuntimeGetFuncPtr54.File = diFileFuncGo55

	// 	!56 = !DISubroutineType(types: !57)
	diSubroutine56 := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramRuntimeGetFuncPtr54.Type = diSubroutine56

	// 	!57 = !{!58, !65}
	diTuple57 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubroutine56.Types = diTuple57

	// 	!58 = !DIDerivedType(tag: DW_TAG_typedef, name: "runtime.funcValue", baseType: !59)
	diDerived58 := &metadata.DIDerivedType{
		MetadataID: -1,
		Tag:        enum.DwarfTagTypedef,
		Name:       "runtime.funcValue",
	}
	diTuple57.Fields = append(diTuple57.Fields, diDerived58)

	// 	!59 = !DICompositeType(tag: DW_TAG_structure_type, size: 64, align: 32, elements: !60)
	diComposite59 := &metadata.DICompositeType{
		MetadataID: -1,
		Tag:        enum.DwarfTagStructureType,
		Size:       64,
		Align:      32,
	}
	diDerived58.BaseType = diComposite59

	// 	!60 = !{!61, !63}
	diTuple60 := &metadata.Tuple{
		MetadataID: -1,
	}
	diComposite59.Elements = diTuple60

	// 	!61 = !DIDerivedType(tag: DW_TAG_member, name: "context", baseType: !62, size: 32, align: 32)
	diDerived61 := &metadata.DIDerivedType{
		MetadataID: -1,
		Tag:        enum.DwarfTagMember,
		Name:       "context",
		Size:       32,
		Align:      32,
	}
	diTuple60.Fields = append(diTuple60.Fields, diDerived61)

	// 	!62 = !DIDerivedType(tag: DW_TAG_pointer_type, name: "unsafe.Pointer", baseType: null, size: 32, align: 32, dwarfAddressSpace: 0)
	diDerived62 := &metadata.DIDerivedType{
		MetadataID:        -1,
		Tag:               enum.DwarfTagPointerType,
		Name:              "unsafe.Pointer",
		Size:              32,
		Align:             32,
		DwarfAddressSpace: 0,
		BaseType:          metadata.Null,
	}
	diDerived61.BaseType = diDerived62

	// 	!63 = !DIDerivedType(tag: DW_TAG_member, name: "id", baseType: !64, size: 32, align: 32, offset: 32)
	diDerived63 := &metadata.DIDerivedType{
		MetadataID: -1,
		Tag:        enum.DwarfTagMember,
		Name:       "id",
		Size:       32,
		Align:      32,
		Offset:     32,
	}
	diTuple60.Fields = append(diTuple60.Fields, diDerived63)

	// 	!64 = !DIBasicType(name: "uintptr", size: 32, encoding: DW_ATE_unsigned)
	diBasicType64 := &metadata.DIBasicType{
		MetadataID: -1,
		Name:       "uintptr",
		Size:       32,
		Encoding:   enum.DwarfAttEncodingUnsigned,
	}
	diDerived63.BaseType = diBasicType64

	// 	!65 = !DIDerivedType(tag: DW_TAG_pointer_type, baseType: !19, size: 32, align: 32, dwarfAddressSpace: 0)
	diDerived65 := &metadata.DIDerivedType{
		MetadataID:        -1,
		Tag:               enum.DwarfTagPointerType,
		BaseType:          diBasicType19,
		Size:              32,
		Align:             32,
		DwarfAddressSpace: 0,
	}
	diTuple57.Fields = append(diTuple57.Fields, diDerived65)

	// 	!66 = !{!67, !68}
	diTuple66 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubprogramRuntimeGetFuncPtr54.RetainedNodes = diTuple66

	// 	!67 = !DILocalVariable(name: "val", arg: 1, scope: !54, file: !55, line: 26, type: !58)
	diLocalVarVal67 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "val",
		Arg:        1,
		Scope:      diSubprogramRuntimeGetFuncPtr54,
		File:       diFileFuncGo55,
		Line:       26,
		Type:       diDerived58,
	}
	dbgLoc["LocalVarVal"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarVal67))
	diTuple66.Fields = append(diTuple66.Fields, diLocalVarVal67)

	// 	!68 = !DILocalVariable(name: "signature", arg: 2, scope: !54, file: !55, line: 26, type: !65)
	diLocalVarSignature68 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "signature",
		Arg:        2,
		Scope:      diSubprogramRuntimeGetFuncPtr54,
		File:       diFileFuncGo55,
		Line:       26,
		Type:       diDerived65,
	}
	dbgLoc["LocalVarSignature"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarSignature68))
	diTuple66.Fields = append(diTuple66.Fields, diLocalVarSignature68)

	// 	!69 = !DILocation(line: 26, column: 6, scope: !54)
	dbgLoc["69"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       26,
		Column:     6,
		Scope:      diSubprogramRuntimeGetFuncPtr54,
	}

	// 	!70 = !DILocation(line: 27, column: 59, scope: !54)
	dbgLoc["70"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       27,
		Column:     59,
		Scope:      diSubprogramRuntimeGetFuncPtr54,
	}

	// 	!71 = distinct !DISubprogram(name: "runtime.memset", linkageName: "memset", scope: !6, file: !6, line: 73, type: !72, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !75)
	diSubprogramRuntimeMemset71 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.memset",
		LinkageName: "memset",
		Scope:       diFileRuntimeWasmGo6,
		File:        diFileRuntimeWasmGo6,
		Line:        73,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.memset"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeMemset71))

	// 	!72 = !DISubroutineType(types: !73)
	diSubroutine72 := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramRuntimeMemset71.Type = diSubroutine72

	// 	!73 = !{!62, !74, !64}
	diTuple73 := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{diDerived62},
	}
	diSubroutine72.Types = diTuple73

	// 	!74 = !DIBasicType(name: "byte", size: 8, encoding: DW_ATE_unsigned)
	diBasicType74 := &metadata.DIBasicType{
		MetadataID: -1,
		Name:       "byte",
		Size:       8,
		Encoding:   enum.DwarfAttEncodingUnsigned,
	}
	diTuple73.Fields = append(diTuple73.Fields, diBasicType74, diBasicType64)

	// 	!75 = !{!76, !77, !78}
	diTuple75 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubprogramRuntimeMemset71.RetainedNodes = diTuple75

	// 	!76 = !DILocalVariable(name: "ptr", arg: 1, scope: !71, file: !6, line: 73, type: !62)
	diLocalVarPtr76 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "ptr",
		Arg:        1,
		Scope:      diSubprogramRuntimeMemset71,
		File:       diFileRuntimeWasmGo6,
		Line:       73,
		Type:       diDerived62,
	}
	dbgLoc["ptr"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarPtr76))
	diTuple75.Fields = append(diTuple75.Fields, diLocalVarPtr76)

	// 	!77 = !DILocalVariable(name: "c", arg: 2, scope: !71, file: !6, line: 73, type: !74)
	diLocalVarC77 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "c",
		Arg:        2,
		Scope:      diSubprogramRuntimeMemset71,
		File:       diFileRuntimeWasmGo6,
		Line:       73,
		Type:       diBasicType74,
	}
	dbgLoc["c"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarC77))
	diTuple75.Fields = append(diTuple75.Fields, diLocalVarC77)

	// 	!78 = !DILocalVariable(name: "size", arg: 3, scope: !71, file: !6, line: 73, type: !64)
	diLocalVarSize78 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "size",
		Arg:        3,
		Scope:      diSubprogramRuntimeMemset71,
		File:       diFileRuntimeWasmGo6,
		Line:       73,
		Type:       diBasicType64,
	}
	dbgLoc["size"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarSize78))
	diTuple75.Fields = append(diTuple75.Fields, diLocalVarSize78)

	// 	!79 = !DILocation(line: 73, column: 6, scope: !71)
	dbgLoc["79"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       73,
		Column:     6,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!80 = !DILocation(line: 0, scope: !71)
	dbgLoc["80"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!81 = !DILocation(line: 74, column: 6, scope: !71)
	dbgLoc["81"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       74,
		Column:     6,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!82 = !DILocation(line: 74, column: 25, scope: !71)
	dbgLoc["82"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       74,
		Column:     25,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!83 = !DILocation(line: 75, column: 26, scope: !71)
	dbgLoc["83"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       75,
		Column:     26,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!84 = !DILocation(line: 75, column: 3, scope: !71)
	dbgLoc["84"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       75,
		Column:     3,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!85 = !DILocation(line: 77, column: 2, scope: !71)
	dbgLoc["85"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       77,
		Column:     2,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!86 = !DILocation(line: 74, column: 33, scope: !71)
	dbgLoc["86"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       74,
		Column:     33,
		Scope:      diSubprogramRuntimeMemset71,
	}

	// 	!87 = distinct !DISubprogram(name: "runtime.nilPanic", linkageName: "runtime.nilPanic", scope: !88, file: !88, line: 39, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramRuntimeNilPanic87 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime.nilPanic",
		LinkageName:   "runtime.nilPanic",
		Line:          39,
		Type:          diSubroutineType7,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["runtime.nilPanic"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeNilPanic87))

	// 	!88 = !DIFile(filename: "panic.go", directory: "../../../runtime")
	diFilePanicGo88 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "panic.go",
		Directory:  "../../../runtime",
	}
	diSubprogramRuntimeNilPanic87.Scope = diFilePanicGo88
	diSubprogramRuntimeNilPanic87.File = diFilePanicGo88

	// 	!89 = !DILocation(line: 40, column: 14, scope: !87)
	dbgLoc["89"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       40,
		Column:     14,
		Scope:      diSubprogramRuntimeNilPanic87,
	}

	// 	!90 = distinct !DISubprogram(name: "runtime.printnl", linkageName: "runtime.printnl", scope: !91, file: !91, line: 199, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramRuntimePrintnl90 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime.printnl",
		LinkageName:   "runtime.printnl",
		Line:          199,
		Type:          diSubroutineType7,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["runtime.printnl"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimePrintnl90))

	// 	!91 = !DIFile(filename: "print.go", directory: "../../../runtime")
	diFilePrint91 := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "print.go",
		Directory:  "../../../runtime",
	}
	diSubprogramRuntimePrintnl90.Scope = diFilePrint91
	diSubprogramRuntimePrintnl90.File = diFilePrint91

	// 	!92 = !DILocation(line: 200, column: 9, scope: !90)
	dbgLoc["92"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       200,
		Column:     9,
		Scope:      diSubprogramRuntimePrintnl90,
	}

	// 	!93 = !DILocation(line: 201, column: 9, scope: !90)
	dbgLoc["93"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       201,
		Column:     9,
		Scope:      diSubprogramRuntimePrintnl90,
	}

	// 	!94 = !DILocation(line: 0, scope: !90)
	dbgLoc["94"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimePrintnl90,
	}

	// 	!95 = distinct !DISubprogram(name: "runtime.printstring", linkageName: "runtime.printstring", scope: !91, file: !91, line: 12, type: !96, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !102)
	diSubprogramRuntimePrintString95 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.printstring",
		LinkageName: "runtime.printstring",
		Scope:       diFilePrint91,
		File:        diFilePrint91,
		Line:        12,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.printstring"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimePrintString95))

	// 	!96 = !DISubroutineType(types: !97)
	diSubroutine96 := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramRuntimePrintString95.Type = diSubroutine96

	// 	!97 = !{!98}
	diTuple97 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubroutine96.Types = diTuple97

	// 	!98 = !DICompositeType(tag: DW_TAG_structure_type, name: "string", size: 64, align: 32, elements: !99)
	diComposite98 := &metadata.DICompositeType{
		MetadataID: -1,
		Tag:        enum.DwarfTagStructureType,
		Name:       "string",
		Size:       64,
		Align:      32,
	}
	diTuple97.Fields = append(diTuple97.Fields, diComposite98)

	// 	!99 = !{!100, !101}
	diTuple99 := &metadata.Tuple{
		MetadataID: -1,
	}
	diComposite98.Elements = diTuple99

	// 	!100 = !DIDerivedType(tag: DW_TAG_member, name: "ptr", baseType: !65, size: 32, align: 32)
	diDerived100 := &metadata.DIDerivedType{
		MetadataID: -1,
		Tag:        enum.DwarfTagMember,
		Name:       "ptr",
		BaseType:   diDerived65,
		Size:       32,
		Align:      32,
	}
	diTuple99.Fields = append(diTuple99.Fields, diDerived100)

	// 	!101 = !DIDerivedType(tag: DW_TAG_member, name: "len", baseType: !64, size: 32, align: 32, offset: 32)
	diDerived101 := &metadata.DIDerivedType{
		MetadataID: -1,
		Tag:        enum.DwarfTagMember,
		Name:       "len",
		BaseType:   diBasicType64,
		Size:       32,
		Align:      32,
		Offset:     32,
	}
	diTuple99.Fields = append(diTuple99.Fields, diDerived101)

	// 	!102 = !{!103}
	diTuple102 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubprogramRuntimePrintString95.RetainedNodes = diTuple102

	// 	!103 = !DILocalVariable(name: "s", arg: 1, scope: !95, file: !91, line: 12, type: !98)
	diLocalVarS103 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "s",
		Arg:        1,
		Scope:      diSubprogramRuntimePrintString95,
		File:       diFilePrint91,
		Line:       12,
		Type:       diComposite98,
	}
	dbgLoc["s"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarS103))
	diTuple102.Fields = append(diTuple102.Fields, diLocalVarS103)

	// 	!104 = !DILocation(line: 12, column: 6, scope: !95)
	dbgLoc["104"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       12,
		Column:     6,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!105 = !DILocation(line: 0, scope: !95)
	dbgLoc["105"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!106 = !DILocation(line: 13, column: 6, scope: !95)
	dbgLoc["106"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       13,
		Column:     6,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!107 = !DILocation(line: 13, column: 16, scope: !95)
	dbgLoc["107"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       13,
		Column:     16,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!108 = !DILocation(line: 14, column: 12, scope: !95)
	dbgLoc["108"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       14,
		Column:     12,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!109 = !DILocation(line: 14, column: 10, scope: !95)
	dbgLoc["109"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       14,
		Column:     10,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!110 = !DILocation(line: 13, column: 26, scope: !95)
	dbgLoc["110"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       13,
		Column:     26,
		Scope:      diSubprogramRuntimePrintString95,
	}

	// 	!111 = distinct !DISubprogram(name: "runtime.putchar", linkageName: "runtime.putchar", scope: !6, file: !6, line: 36, type: !112, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !114)
	diSubprogramRuntimePutchar111 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.putchar",
		LinkageName: "runtime.putchar",
		Scope:       diFileRuntimeWasmGo6,
		File:        diFileRuntimeWasmGo6,
		Line:        36,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.putchar"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimePutchar111))

	// 	!112 = !DISubroutineType(types: !113)
	diSubroutine112 := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramRuntimePutchar111.Type = diSubroutine112

	// 	!113 = !{!74}
	diTuple113 := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{diBasicType74},
	}
	diSubroutine112.Types = diTuple113

	// 	!114 = !{!115}
	diTuple114 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubprogramRuntimePutchar111.RetainedNodes = diTuple114

	// 	!115 = !DILocalVariable(name: "c", arg: 1, scope: !111, file: !6, line: 36, type: !74)
	diLocalVarC115 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "c",
		Arg:        1,
		Scope:      diSubprogramRuntimePutchar111,
		File:       diFileRuntimeWasmGo6,
		Line:       36,
		Type:       diBasicType74,
	}
	dbgLoc["c115"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarC115))
	diTuple114.Fields = append(diTuple114.Fields, diLocalVarC115)

	// 	!116 = !DILocation(line: 36, column: 6, scope: !111)
	dbgLoc["116"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       36,
		Column:     6,
		Scope:      diSubprogramRuntimePutchar111,
	}

	// 	!117 = !DILocation(line: 0, scope: !111)
	dbgLoc["117"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimePutchar111,
	}

	// 	!118 = !DILocation(line: 37, column: 17, scope: !111)
	dbgLoc["118"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       37,
		Column:     17,
		Scope:      diSubprogramRuntimePutchar111,
	}

	// 	!119 = !DILocation(line: 37, column: 16, scope: !111)
	dbgLoc["119"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       37,
		Column:     16,
		Scope:      diSubprogramRuntimePutchar111,
	}

	// 	!120 = distinct !DISubprogram(name: "runtime.resume", linkageName: "resume", scope: !6, file: !6, line: 48, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramRuntimeResume120 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime.resume",
		LinkageName:   "resume",
		Scope:         diFileRuntimeWasmGo6,
		File:          diFileRuntimeWasmGo6,
		Line:          48,
		Type:          diSubroutineType7,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["runtime.resume"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeResume120))

	// 	!121 = !DILocation(line: 49, column: 13, scope: !120)
	dbgLoc["121"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       49,
		Column:     13,
		Scope:      diSubprogramRuntimeResume120,
	}

	// 	!122 = distinct !DISubprogram(name: "runtime.runtimePanic", linkageName: "runtime.runtimePanic", scope: !88, file: !88, line: 17, type: !96, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !123)
	diSubprogramRuntimeRuntimePanic122 := &metadata.DISubprogram{
		MetadataID:  -1,
		Distinct:    true,
		Name:        "runtime.runtimePanic",
		LinkageName: "runtime.runtimePanic",
		Scope:       diFilePanicGo88,
		File:        diFilePanicGo88,
		Line:        17,
		Type:        diSubroutine96,
		Flags:       enum.DIFlagPrototyped,
		SPFlags:     enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:        diCompileUnit,
	}
	dbgLoc["runtime.runtimePanic"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeRuntimePanic122))

	// 	!123 = !{!124}
	diTuple123 := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubprogramRuntimeRuntimePanic122.RetainedNodes = diTuple123

	// 	!124 = !DILocalVariable(name: "msg", arg: 1, scope: !122, file: !88, line: 17, type: !98)
	diLocalVarMsg124 := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "msg",
		Arg:        1,
		Scope:      diSubprogramRuntimeRuntimePanic122,
		File:       diFilePanicGo88,
		Line:       17,
		Type:       diComposite98,
	}
	dbgLoc["msg"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarMsg124))
	diTuple123.Fields = append(diTuple123.Fields, diLocalVarMsg124)

	// 	!125 = !DILocation(line: 17, column: 6, scope: !122)
	dbgLoc["125"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       17,
		Column:     6,
		Scope:      diSubprogramRuntimeRuntimePanic122,
	}

	// 	!126 = !DILocation(line: 18, column: 13, scope: !122)
	dbgLoc["126"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       18,
		Column:     13,
		Scope:      diSubprogramRuntimeRuntimePanic122,
	}

	// 	!127 = !DILocation(line: 19, column: 9, scope: !122)
	dbgLoc["127"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       19,
		Column:     9,
		Scope:      diSubprogramRuntimeRuntimePanic122,
	}

	// 	!128 = !DILocation(line: 69, column: 6, scope: !129, inlinedAt: !130)
	dbgLoc["128"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       69,
		Column:     6,
	}

	// 	!129 = distinct !DISubprogram(name: "runtime.abort", linkageName: "runtime.abort", scope: !6, file: !6, line: 68, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
	diSubprogramRuntimeAbort129 := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "runtime.abort",
		LinkageName:   "runtime.abort",
		Scope:         diFileRuntimeWasmGo6,
		File:          diFileRuntimeWasmGo6,
		Line:          68,
		Type:          diSubroutineType7,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagLocalToUnit | enum.DISPFlagDefinition | enum.DISPFlagOptimized,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple2,
	}
	dbgLoc["runtime.abort"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramRuntimeAbort129))
	dbgLoc["128"].Scope = diSubprogramRuntimeAbort129

	// 	!130 = distinct !DILocation(line: 20, column: 7, scope: !122)
	dbgLoc["130"] = &metadata.DILocation{
		MetadataID: -1,
		Distinct:   true,
		Line:       20,
		Column:     7,
		Scope:      diSubprogramRuntimeRuntimePanic122,
	}
	dbgLoc["128"].InlinedAt = dbgLoc["130"]

	// 	!131 = !DILocation(line: 0, scope: !129, inlinedAt: !130)
	// TODO: Create a PR to fix the bug where "Line: 0" gets omitted
	dbgLoc["131"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       0,
		LineValid:  true,
		Scope:      diSubprogramRuntimeAbort129,
		InlinedAt:  dbgLoc["130"],
	}

	// * Named metadata definitions *

	// 	!llvm.dbg.cu = !{!0}
	llvmDbgCu := &metadata.NamedDef{
		Name:  "llvm.dbg.cu",
		Nodes: []metadata.Node{diCompileUnit},
	}
	m.NamedMetadataDefs["llvm.dbg.cu"] = llvmDbgCu

	// 	!llvm.module.flags = !{!3, !4}
	llvmModuleFlags := &metadata.NamedDef{
		Name:  "llvm.module.flags",
		Nodes: []metadata.Node{debugInfoVersion, dwarfVersion},
	}
	m.NamedMetadataDefs["llvm.module.flags"] = llvmModuleFlags

	// Attach the debugging info to the module
	m.MetadataDefs = append(m.MetadataDefs, diCompileUnit)
	m.MetadataDefs = append(m.MetadataDefs, diFileMainGo1)
	m.MetadataDefs = append(m.MetadataDefs, emptyTuple2)
	m.MetadataDefs = append(m.MetadataDefs, debugInfoVersion)
	m.MetadataDefs = append(m.MetadataDefs, dwarfVersion)
	m.MetadataDefs = append(m.MetadataDefs, diSubprogram_start)
	m.MetadataDefs = append(m.MetadataDefs, diFileRuntimeWasmGo6)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutineType7)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["8"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeInitAll9)
	m.MetadataDefs = append(m.MetadataDefs, diFileRuntimeGo10)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["11"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["12"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeActivateTask13)
	m.MetadataDefs = append(m.MetadataDefs, diFileSchedulerGo14)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutine15)
	m.MetadataDefs = append(m.MetadataDefs, diTuple16)
	m.MetadataDefs = append(m.MetadataDefs, diDerived17)
	m.MetadataDefs = append(m.MetadataDefs, diDerived18)
	m.MetadataDefs = append(m.MetadataDefs, diBasicType19)
	m.MetadataDefs = append(m.MetadataDefs, diTuple20)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarTask21)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["22"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["23"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["24"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["25"])
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarT26)
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeRunqueuePushBack27)
	m.MetadataDefs = append(m.MetadataDefs, diTuple28)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["29"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["30"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["31"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["32"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["33"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["34"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["35"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["36"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["37"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["38"])
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarT39)
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeCoroutinePromisePtr40)
	m.MetadataDefs = append(m.MetadataDefs, diTuple41)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["42"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["43"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["44"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["45"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["46"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeCwaMain47)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["48"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["49"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["50"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramMainGoMain51)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["52"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["53"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeGetFuncPtr54)
	m.MetadataDefs = append(m.MetadataDefs, diFileFuncGo55)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutine56)
	m.MetadataDefs = append(m.MetadataDefs, diTuple57)
	m.MetadataDefs = append(m.MetadataDefs, diDerived58)
	m.MetadataDefs = append(m.MetadataDefs, diComposite59)
	m.MetadataDefs = append(m.MetadataDefs, diTuple60)
	m.MetadataDefs = append(m.MetadataDefs, diDerived61)
	m.MetadataDefs = append(m.MetadataDefs, diDerived62)
	m.MetadataDefs = append(m.MetadataDefs, diDerived63)
	m.MetadataDefs = append(m.MetadataDefs, diBasicType64)
	m.MetadataDefs = append(m.MetadataDefs, diDerived65)
	m.MetadataDefs = append(m.MetadataDefs, diTuple66)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarVal67)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarSignature68)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["69"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["70"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeMemset71)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutine72)
	m.MetadataDefs = append(m.MetadataDefs, diTuple73)
	m.MetadataDefs = append(m.MetadataDefs, diBasicType74)
	m.MetadataDefs = append(m.MetadataDefs, diTuple75)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarPtr76)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarC77)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarSize78)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["79"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["80"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["81"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["82"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["83"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["84"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["85"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["86"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeNilPanic87)
	m.MetadataDefs = append(m.MetadataDefs, diFilePanicGo88)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["89"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimePrintnl90)
	m.MetadataDefs = append(m.MetadataDefs, diFilePrint91)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["92"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["93"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["94"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimePrintString95)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutine96)
	m.MetadataDefs = append(m.MetadataDefs, diTuple97)
	m.MetadataDefs = append(m.MetadataDefs, diComposite98)
	m.MetadataDefs = append(m.MetadataDefs, diTuple99)
	m.MetadataDefs = append(m.MetadataDefs, diDerived100)
	m.MetadataDefs = append(m.MetadataDefs, diDerived101)
	m.MetadataDefs = append(m.MetadataDefs, diTuple102)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarS103)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["104"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["105"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["106"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["107"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["108"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["109"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["110"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimePutchar111)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutine112)
	m.MetadataDefs = append(m.MetadataDefs, diTuple113)
	m.MetadataDefs = append(m.MetadataDefs, diTuple114)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarC115)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["116"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["117"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["118"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["119"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeResume120)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["121"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeRuntimePanic122)
	m.MetadataDefs = append(m.MetadataDefs, diTuple123)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarMsg124)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["125"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["126"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["127"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["128"])
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramRuntimeAbort129)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["130"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["131"])
	return
}

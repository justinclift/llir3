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
	// Just for convenience
	i32 = types.I32

	// Keep track of debug locations
	dbgLoc = make(map[string]*metadata.DILocation)
)

func main() {
	// Convenience constants.
	var (
		zero = constant.NewInt(i32, 0)
		twelve = constant.NewInt(i32, 12)
		thirty = constant.NewInt(i32, 30)
	)

	// Create LLVM IR module.
	m := ir.NewModule()

	// Add metadata
	addMetadata(m)
	m.DataLayout = "e-m:e-i64:64-f80:128-n8:16:32:64-S128"
	m.TargetTriple = "x86_64-pc-linux-gnu"
	m.SourceFilename = "target.c"

	// Empty DIExpression
	//    !DIExpression()
	emptyExpr := &metadata.DIExpression{
		MetadataID: -1,
	}

	// Define a void function, used as the return type for calls further down
	// declare void @llvm.dbg.declare(metadata, metadata, metadata) #1
	llvmDbgDeclare := m.NewFunc(
		"llvm.dbg.declare",
		types.Void,
		ir.NewParam("", types.Metadata),
		ir.NewParam("", types.Metadata),
		ir.NewParam("", types.Metadata),
	)

	// Add the function attributes
	// ; Function Attrs: nounwind readnone speculatable
	llvmDbgDeclare.FuncAttrs = append(llvmDbgDeclare.FuncAttrs, enum.FuncAttrNoUnwind)
	llvmDbgDeclare.FuncAttrs = append(llvmDbgDeclare.FuncAttrs, enum.FuncAttrReadNone)
	llvmDbgDeclare.FuncAttrs = append(llvmDbgDeclare.FuncAttrs, enum.FuncAttrSpeculatable)

	// Define the "foo" function
	//   `int foo(int a, int b)`
	// define dso_local i32 @foo(i32, i32) #0 !dbg !7 {
	aParam := ir.NewParam("", i32)
	bParam := ir.NewParam("", i32)

	// define dso_local i32 @foo(i32, i32) #0 !dbg !7 {
	fooFunc := m.NewFunc("foo", i32, aParam, bParam)
	fooFunc.Metadata = append(fooFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["SPFoo"]})

	// Add the function attributes
	// ; Function Attrs: noinline nounwind optnone uwtable
	fooFunc.FuncAttrs = append(fooFunc.FuncAttrs, enum.FuncAttrNoInline)
	fooFunc.FuncAttrs = append(fooFunc.FuncAttrs, enum.FuncAttrNoUnwind)
	fooFunc.FuncAttrs = append(fooFunc.FuncAttrs, enum.FuncAttrOptNone)
	fooFunc.FuncAttrs = append(fooFunc.FuncAttrs, enum.FuncAttrSpeculatable)

	// Create LLVM Block, for containing the subsequent function code
	fooEntry := fooFunc.NewBlock("")

	// %3 = alloca i32, align 4
	a := fooEntry.NewAlloca(i32)
	a.Align = 4

	// %4 = alloca i32, align 4
	b := fooEntry.NewAlloca(i32)
	b.Align = 4

	// %5 = alloca i32, align 4
	sum := fooEntry.NewAlloca(i32)
	sum.Align = 4

	// store i32 %0, i32* %3, align 4
	z1 := fooEntry.NewStore(aParam, a)
	z1.Align = 4

	// call void @llvm.dbg.declare(metadata i32* %3, metadata !11, metadata !DIExpression()), !dbg !12
	dbgDeclareA := fooEntry.NewCall(llvmDbgDeclare, &metadata.Value{Value: a}, &metadata.Value{Value: dbgLoc["LocalVarA"]}, &metadata.Value{Value: emptyExpr})
	dbgDeclareA.Metadata = append(dbgDeclareA.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["A"]})

	// store i32 %1, i32* %4, align 4
	z2 := fooEntry.NewStore(bParam, b)
	z2.Align = 4

	// call void @llvm.dbg.declare(metadata i32* %4, metadata !13, metadata !DIExpression()), !dbg !14
	dbgDeclareB := fooEntry.NewCall(llvmDbgDeclare, &metadata.Value{Value: b}, &metadata.Value{Value: dbgLoc["LocalVarB"]}, &metadata.Value{Value: emptyExpr})
	dbgDeclareB.Metadata = append(dbgDeclareB.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["B"]})

	// call void @llvm.dbg.declare(metadata i32* %5, metadata !15, metadata !DIExpression()), !dbg !16
	dbgDeclareSum := fooEntry.NewCall(llvmDbgDeclare, &metadata.Value{Value: sum}, &metadata.Value{Value: dbgLoc["LocalVarSum"]}, &metadata.Value{Value: emptyExpr})
	dbgDeclareSum.Metadata = append(dbgDeclareSum.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Sum"]})

	// %6 = load i32, i32* %3, align 4, !dbg !17
	aVal := fooEntry.NewLoad(a)
	aVal.Align = 4
	aVal.Metadata = append(aVal.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["A2"]})

	// %7 = load i32, i32* %4, align 4, !dbg !18
	bVal := fooEntry.NewLoad(b)
	bVal.Align = 4
	bVal.Metadata = append(bVal.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["B2"]})

	// %8 = add nsw i32 %6, %7, !dbg !19
	tmp := fooEntry.NewAdd(aVal, bVal)
	tmp.OverflowFlags = append(tmp.OverflowFlags, enum.OverflowFlagNSW)
	tmp.Metadata = append(tmp.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Plus"]})

	// store i32 %8, i32* %5, align 4, !dbg !20
	z3 := fooEntry.NewStore(tmp, sum)
	z3.Align = 4
	z3.Metadata = append(z3.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Equals"]})

	// %9 = load i32, i32* %5, align 4, !dbg !21
	sumVal := fooEntry.NewLoad(sum)
	sumVal.Align = 4
	sumVal.Metadata = append(sumVal.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Sum2"]})

	// ret i32 %9, !dbg !22
	z4 := fooEntry.NewRet(sumVal)
	z4.Metadata = append(z4.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Return"]})

	// Define the "main" function
	//   `int main()`
	// define dso_local i32 @main() #0 !dbg !23 {
	mainFunc := m.NewFunc("main", i32)
	mainFunc.Metadata = append(mainFunc.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["SPMain"]})

	// Add the function attributes
	// ; Function Attrs: noinline nounwind optnone uwtable
	mainFunc.FuncAttrs = append(mainFunc.FuncAttrs, enum.FuncAttrNoInline)
	mainFunc.FuncAttrs = append(mainFunc.FuncAttrs, enum.FuncAttrNoUnwind)
	mainFunc.FuncAttrs = append(mainFunc.FuncAttrs, enum.FuncAttrOptNone)
	mainFunc.FuncAttrs = append(mainFunc.FuncAttrs, enum.FuncAttrSpeculatable)

	// Create LLVM Block, for containing the subsequent function code
	mainEntry := mainFunc.NewBlock("")

	// %1 = alloca i32, align 4
	a2 := mainEntry.NewAlloca(i32)
	a2.Align = 4

	// store i32 %0, i32* %3, align 4
	s2 := mainEntry.NewStore(zero, a2)
	s2.Align = 4

	// %2 = call i32 @foo(i32 12, i32 30), !dbg !26
	retVal := mainEntry.NewCall(fooFunc, twelve, thirty)
	retVal.Metadata = append(retVal.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Foo"]})

	// ret i32 %2, !dbg !27
	z5 := mainEntry.NewRet(retVal)
	z5.Metadata = append(z5.Metadata, &metadata.Attachment{Name: "dbg", Node: dbgLoc["Return2"]})

	// Print LLVM IR assembly to standard output.
	fmt.Println(m)
}

func addMetadata(m *ir.Module) {
	// Note, the reason we specify MetadataID to be -1, is so that the IR
	// package may assign the metadata definition an arbitrary unique ID (and 0
	// is a valid ID).
	//
	// I think we should try to find a cleaner way to handle this. Any
	// suggestions are warmly welcome! :)

	// Convenience constants.
	var (
		one   = constant.NewInt(i32, 1)
		two   = constant.NewInt(i32, 2)
		three = constant.NewInt(i32, 3)
		four  = constant.NewInt(i32, 4)
	)

	// Unnamed metadata definitions.

	// DICompileUnit
	// !0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git
	//   295e7b4486abfc0f01e6d22276931165f324c61e)", isOptimized: false, runtimeVersion: 0, emissionKind: FullDebug, enums: !2, nameTableKind: None)
	diCompileUnit := &metadata.DICompileUnit{
		MetadataID:     -1,
		Distinct:       true,
		Language:       enum.DwarfLangC99,
		Producer:       "clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git 295e7b4486abfc0f01e6d22276931165f324c61e)",
		IsOptimized:    false,
		RuntimeVersion: 0,
		EmissionKind:   enum.EmissionKindFullDebug,
		NameTableKind:  enum.NameTableKindNone,
	}

	// DIFile
	// !1 = !DIFile(filename: "target.c", directory: "/home/jc/git_repos2/llir3")
	diFile := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "target.c",
		Directory:  "/home/jc/git_repos2/llir3",
	}
	diCompileUnit.File = diFile

	// Empty tuple
	// !2 = !{}
	emptyTuple := &metadata.Tuple{
		MetadataID: -1,
	}
	diCompileUnit.Enums = emptyTuple

	// DWARF metadata
	// !3 = !{i32 2, !"Dwarf Version", i32 4}
	dwarfVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{two, &metadata.String{Value: "Dwarf Version"}, four},
	}
	// !4 = !{i32 2, !"Debug Info Version", i32 3}
	debugInfoVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{two, &metadata.String{Value: "Debug Info Version"}, three},
	}
	// !5 = !{i32 1, !"wchar_size", i32 4}
	wcharSize := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{one, &metadata.String{Value: "wchar_size"}, four},
	}
	// !6 = !{!"clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git 295e7b4486abfc0f01e6d22276931165f324c61e)"}
	clangVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{&metadata.String{Value: "clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git 295e7b4486abfc0f01e6d22276931165f324c61e)"}},
	}

	// DISubprogram
	// !7 = distinct !DISubprogram(name: "foo", scope: !1, file: !1, line: 1, type: !8, scopeLine: 1, flags: DIFlagPrototyped, spFlags: DISPFlagDefinition, unit: !0, retainedNodes: !2)
	diSubprogramFoo := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "foo",
		Scope:         diFile,
		File:          diFile,
		Line:          1,
		ScopeLine:     1,
		Flags:         enum.DIFlagPrototyped,
		SPFlags:       enum.DISPFlagDefinition,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple,
	}
	dbgLoc["SPFoo"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramFoo))

	// DISubroutineType
	// !8 = !DISubroutineType(types: !9)
	diSubroutineType := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramFoo.Type = diSubroutineType

	// Types tuple.
	// !9 = !{!10, !10, !10}
	typesTuple := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubroutineType.Types = typesTuple

	// DIBasicType
	// !10 = !DIBasicType(name: "int", size: 32, encoding: DW_ATE_signed)
	diBasicTypeI32 := &metadata.DIBasicType{
		MetadataID: -1,
		Name:       "int",
		Size:       32,
		Encoding:   enum.DwarfAttEncodingSigned,
	}
	typesTuple.Fields = []metadata.Field{diBasicTypeI32, diBasicTypeI32, diBasicTypeI32}

	// DILocalVariable
	// !11 = !DILocalVariable(name: "a", arg: 1, scope: !7, file: !1, line: 1, type: !10)
	diLocalVarA := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "a",
		Arg:        1,
		Scope:      diSubprogramFoo,
		File:       diFile,
		Line:       1,
		Type:       diBasicTypeI32,
	}
	dbgLoc["LocalVarA"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarA))

	// DILocation
	// !12 = !DILocation(line: 1, column: 13, scope: !7)
	dbgLoc["A"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       1,
		Column:     13,
		Scope:      diSubprogramFoo,
	}

	// DILocalVariable
	// !13 = !DILocalVariable(name: "b", arg: 2, scope: !7, file: !1, line: 1, type: !10)
	diLocalVarB := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "b",
		Arg:        2,
		Scope:      diSubprogramFoo,
		File:       diFile,
		Line:       1,
		Type:       diBasicTypeI32,
	}
	dbgLoc["LocalVarB"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarB))

	// DILocation
	// !14 = !DILocation(line: 1, column: 20, scope: !7)
	dbgLoc["B"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       1,
		Column:     20,
		Scope:      diSubprogramFoo,
	}

	// DILocalVariable
	// !15 = !DILocalVariable(name: "sum", scope: !7, file: !1, line: 2, type: !10)
	diLocalVarSum := &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "sum",
		Scope:      diSubprogramFoo,
		File:       diFile,
		Line:       2,
		Type:       diBasicTypeI32,
	}
	dbgLoc["LocalVarSum"] = (*metadata.DILocation)(unsafe.Pointer(diLocalVarSum))

	// DILocation
	// !16 = !DILocation(line: 2, column: 9, scope: !7)
	diLocSum := &metadata.DILocation{
		MetadataID: -1,
		Line:       2,
		Column:     9,
		Scope:      diSubprogramFoo,
	}
	dbgLoc["Sum"] = (*metadata.DILocation)(unsafe.Pointer(diLocSum))

	// DILocation
	// !17 = !DILocation(line: 3, column: 11, scope: !7)
	dbgLoc["A2"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       3,
		Column:     11,
		Scope:      diSubprogramFoo,
	}

	// DILocation
	// !18 = !DILocation(line: 3, column: 15, scope: !7)
	dbgLoc["B2"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       3,
		Column:     15,
		Scope:      diSubprogramFoo,
	}

	// DILocation
	// !19 = !DILocation(line: 3, column: 13, scope: !7)
	dbgLoc["Plus"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       3,
		Column:     13,
		Scope:      diSubprogramFoo,
	}

	// DILocation
	// !20 = !DILocation(line: 3, column: 9, scope: !7)
	dbgLoc["Equals"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       3,
		Column:     9,
		Scope:      diSubprogramFoo,
	}

	// DILocation
	// !21 = !DILocation(line: 4, column: 12, scope: !7)
	dbgLoc["Sum2"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       4,
		Column:     12,
		Scope:      diSubprogramFoo,
	}

	// DILocation
	// !22 = !DILocation(line: 4, column: 5, scope: !7)
	dbgLoc["Return"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       4,
		Column:     5,
		Scope:      diSubprogramFoo,
	}

	m.MetadataDefs = append(m.MetadataDefs, diCompileUnit)
	m.MetadataDefs = append(m.MetadataDefs, diFile)
	m.MetadataDefs = append(m.MetadataDefs, emptyTuple)
	m.MetadataDefs = append(m.MetadataDefs, dwarfVersion)
	m.MetadataDefs = append(m.MetadataDefs, debugInfoVersion)
	m.MetadataDefs = append(m.MetadataDefs, wcharSize)
	m.MetadataDefs = append(m.MetadataDefs, clangVersion)
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramFoo)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutineType)
	m.MetadataDefs = append(m.MetadataDefs, typesTuple)
	m.MetadataDefs = append(m.MetadataDefs, diBasicTypeI32)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarA)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["A"])
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarB)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["B"])
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarSum)
	m.MetadataDefs = append(m.MetadataDefs, diLocSum)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["A2"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["B2"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["Plus"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["Equals"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["Sum2"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["Return"])

	// DISubprogram
	// !23 = distinct !DISubprogram(name: "main", scope: !1, file: !1, line: 7, type: !24, scopeLine: 7, spFlags: DISPFlagDefinition, unit: !0, retainedNodes: !2)
	diSubprogramMain := &metadata.DISubprogram{
		MetadataID:    -1,
		Distinct:      true,
		Name:          "main",
		Scope:         diFile,
		File:          diFile,
		Line:          7,
		ScopeLine:     7,
		SPFlags:       enum.DISPFlagDefinition,
		Unit:          diCompileUnit,
		RetainedNodes: emptyTuple,
	}
	dbgLoc["SPMain"] = (*metadata.DILocation)(unsafe.Pointer(diSubprogramMain))

	// DISubroutineType
	// !24 = !DISubroutineType(types: !25)
	diSubroutineTypeMain := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramMain.Type = diSubroutineTypeMain

	// Types tuple
	// !25 = !{!10}
	typesTupleMain := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubroutineTypeMain.Types = typesTupleMain
	typesTupleMain.Fields = []metadata.Field{diBasicTypeI32}

	// DILocation
	// !26 = !DILocation(line: 8, column: 12, scope: !23)
	dbgLoc["Foo"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       8,
		Column:     12,
		Scope:      diSubprogramMain,
	}

	// DILocation
	// !27 = !DILocation(line: 8, column: 5, scope: !23)
	dbgLoc["Return2"] = &metadata.DILocation{
		MetadataID: -1,
		Line:       8,
		Column:     5,
		Scope:      diSubprogramMain,
	}

	// Attach the debugging info to the module
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramMain)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutineTypeMain)
	m.MetadataDefs = append(m.MetadataDefs, typesTupleMain)
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["Foo"])
	m.MetadataDefs = append(m.MetadataDefs, dbgLoc["Return2"])

	// * Named metadata definitions *

	// !llvm.dbg.cu = !{!0}
	llvmDbgCu := &metadata.NamedDef{
		Name:  "llvm.dbg.cu",
		Nodes: []metadata.Node{diCompileUnit},
	}
	m.NamedMetadataDefs["llvm.dbg.cu"] = llvmDbgCu

	// !llvm.module.flags = !{!3, !4, !5}
	llvmModuleFlags := &metadata.NamedDef{
		Name:  "llvm.module.flags",
		Nodes: []metadata.Node{dwarfVersion, debugInfoVersion, wcharSize},
	}
	m.NamedMetadataDefs["llvm.module.flags"] = llvmModuleFlags

	// !llvm.ident = !{!6}
	llvmIdent := &metadata.NamedDef{
		Name:  "llvm.ident",
		Nodes: []metadata.Node{clangVersion},
	}
	m.NamedMetadataDefs["llvm.ident"] = llvmIdent

	return
}

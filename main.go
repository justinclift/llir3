package main

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/metadata"
	"github.com/llir/llvm/ir/types"
)

// Convenience types.
var (
	i32 = types.I32
)

func main() {
	// Convenience constants.
	var (
		twelve = constant.NewInt(i32, 12)
		thirty = constant.NewInt(i32, 30)
	)

	// Create LLVM IR module.
	m := ir.NewModule()

	// Add metadata.
	diLocalVarA, diLocalVarB, diLocalVarSum, diLocA, diLocB, diLocSum := addMetadata(m)
	// Empty DIExpression
	//    !DIExpression()
	emptyExpr := &metadata.DIExpression{
		MetadataID: -1,
	}

	// Declare llvm.dbg.declare function.
	//    declare void @llvm.dbg.declare(metadata, metadata, metadata)
	llvmDbgDeclare := m.NewFunc(
		"llvm.dbg.declare",
		types.Void,
		ir.NewParam("", types.Metadata),
		ir.NewParam("", types.Metadata),
		ir.NewParam("", types.Metadata),
	)

	// Define foo function.
	//    int foo(int a, int b)
	// TODO: uncomment.
	//aParam := ir.NewParam("a", i32)
	//bParam := ir.NewParam("b", i32)
	aParam := ir.NewParam("", i32)
	bParam := ir.NewParam("", i32)
	fooFunc := m.NewFunc("foo", i32, aParam, bParam)
	fooEntry := fooFunc.NewBlock("")
	a := fooEntry.NewAlloca(i32)
	b := fooEntry.NewAlloca(i32)
	sum := fooEntry.NewAlloca(i32)
	fooEntry.NewStore(aParam, a)
	dbgDeclareA := fooEntry.NewCall(llvmDbgDeclare, &metadata.Value{Value: a}, &metadata.Value{Value: diLocalVarA}, &metadata.Value{Value: emptyExpr})
	dbgDeclareA.Metadata = append(dbgDeclareA.Metadata, &metadata.Attachment{Name: "dbg", Node: diLocA})
	fooEntry.NewStore(bParam, b)
	dbgDeclareB := fooEntry.NewCall(llvmDbgDeclare, &metadata.Value{Value: b}, &metadata.Value{Value: diLocalVarB}, &metadata.Value{Value: emptyExpr})
	dbgDeclareB.Metadata = append(dbgDeclareB.Metadata, &metadata.Attachment{Name: "dbg", Node: diLocB})
	dbgDeclareSum := fooEntry.NewCall(llvmDbgDeclare, &metadata.Value{Value: sum}, &metadata.Value{Value: diLocalVarSum}, &metadata.Value{Value: emptyExpr})
	dbgDeclareSum.Metadata = append(dbgDeclareSum.Metadata, &metadata.Attachment{Name: "dbg", Node: diLocSum})
	aVal := fooEntry.NewLoad(a)
	bVal := fooEntry.NewLoad(b)
	tmp := fooEntry.NewAdd(aVal, bVal)
	fooEntry.NewStore(tmp, sum)
	sumVal := fooEntry.NewLoad(sum)
	fooEntry.NewRet(sumVal)

	// Define main function
	//    int main()
	mainFunc := m.NewFunc("main", i32)
	mainEntry := mainFunc.NewBlock("")
	retVal := mainEntry.NewCall(fooFunc, twelve, thirty)
	mainEntry.NewRet(retVal)

	// Print LLVM IR assembly to standard output.
	fmt.Println(m)
}

func addMetadata(m *ir.Module) (diLocalVarA, diLocalVarB, diLocalVarSum *metadata.DILocalVariable, diLocA, diLocB, diLocSum *metadata.DILocation) {
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
		seven = constant.NewInt(i32, 7)
	)

	// Unnamed metadata definitions.

	// DICompileUnit
	//    !0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "clang version 8.0.0 (tags/RELEASE_800/final)", isOptimized: false, runtimeVersion: 0, emissionKind: FullDebug, enums: !2, nameTableKind: None)
	diCompileUnit := &metadata.DICompileUnit{
		MetadataID:   -1,
		Distinct:     true,
		Language:     enum.DwarfLangC99,
		Producer:     "clang version 8.0.0 (tags/RELEASE_800/final)",
		EmissionKind: enum.EmissionKindFullDebug,
	}

	// DIFile
	//    !1 = !DIFile(filename: "foo.c", directory: "/home/u/Desktop/foo")
	diFile := &metadata.DIFile{
		MetadataID: -1,
		Filename:   "foo.c",
		Directory:  "/home/u/Desktop/foo",
	}
	diCompileUnit.File = diFile

	// Empty tuple.
	//    !2 = !{}
	emptyTuple := &metadata.Tuple{
		MetadataID: -1,
	}
	diCompileUnit.Enums = emptyTuple

	// Dwarf metadata.
	//    !3 = !{i32 2, !"Dwarf Version", i32 4}
	dwarfVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{two, &metadata.String{Value: "Dwarf Version"}, four},
	}
	//    !4 = !{i32 2, !"Debug Info Version", i32 3}
	debugInfoVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{two, &metadata.String{Value: "Debug Info Version"}, three},
	}
	//    !5 = !{i32 1, !"wchar_size", i32 4}
	wcharSize := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{one, &metadata.String{Value: "wchar_size"}, four},
	}
	//    !6 = !{i32 7, !"PIC Level", i32 2}
	picLevel := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{seven, &metadata.String{Value: "PIC Level"}, two},
	}
	//    !7 = !{i32 7, !"PIE Level", i32 2}
	pieLevel := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{seven, &metadata.String{Value: "PIE Level"}, two},
	}
	//    !8 = !{!"clang version 8.0.0 (tags/RELEASE_800/final)"}
	clangVersion := &metadata.Tuple{
		MetadataID: -1,
		Fields:     []metadata.Field{&metadata.String{Value: "clang version 8.0.0 (tags/RELEASE_800/final)"}},
	}

	// DISubprogram
	//    !9 = distinct !DISubprogram(name: "foo", scope: !1, file: !1, line: 1, type: !10, scopeLine: 1, flags: DIFlagPrototyped, spFlags: DISPFlagDefinition, unit: !0, retainedNodes: !2)
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

	// DISubroutineType
	//    !10 = !DISubroutineType(types: !11)
	diSubroutineType := &metadata.DISubroutineType{
		MetadataID: -1,
	}
	diSubprogramFoo.Type = diSubroutineType

	// Types tuple.
	//    !11 = !{!12, !12, !12}
	typesTuple := &metadata.Tuple{
		MetadataID: -1,
	}
	diSubroutineType.Types = typesTuple

	// DIBasicType
	//    !12 = !DIBasicType(name: "int", size: 32, encoding: DW_ATE_signed)
	diBasicTypeI32 := &metadata.DIBasicType{
		MetadataID: -1,
		Name:       "int",
		Size:       32,
		Encoding:   enum.DwarfAttEncodingSigned,
	}
	typesTuple.Fields = []metadata.Field{diBasicTypeI32, diBasicTypeI32, diBasicTypeI32}

	// DILocalVariable
	//    !13 = !DILocalVariable(name: "a", arg: 1, scope: !9, file: !1, line: 1, type: !12)
	diLocalVarA = &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "a",
		Arg:        1,
		Scope:      diSubprogramFoo,
		File:       diFile,
		Line:       1,
		Type:       diBasicTypeI32,
	}

	// DILocation
	//    !14 = !DILocation(line: 1, column: 13, scope: !9)
	diLocA = &metadata.DILocation{
		MetadataID: -1,
		Line:       1,
		Column:     13,
		Scope:      diSubprogramFoo,
	}

	// DILocalVariable
	//    !15 = !DILocalVariable(name: "b", arg: 2, scope: !9, file: !1, line: 1, type: !12)
	diLocalVarB = &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "b",
		Arg:        2,
		Scope:      diSubprogramFoo,
		File:       diFile,
		Line:       1,
		Type:       diBasicTypeI32,
	}

	// DILocation
	//    !16 = !DILocation(line: 1, column: 20, scope: !9)
	diLocB = &metadata.DILocation{
		MetadataID: -1,
		Line:       1,
		Column:     20,
		Scope:      diSubprogramFoo,
	}

	// DILocalVariable
	//    !17 = !DILocalVariable(name: "sum", scope: !9, file: !1, line: 2, type: !12)
	diLocalVarSum = &metadata.DILocalVariable{
		MetadataID: -1,
		Name:       "sum",
		Scope:      diSubprogramFoo,
		File:       diFile,
		Line:       2,
		Type:       diBasicTypeI32,
	}

	// DILocation
	//    !18 = !DILocation(line: 2, column: 6, scope: !9)
	diLocSum = &metadata.DILocation{
		MetadataID: -1,
		Line:       2,
		Column:     6,
		Scope:      diSubprogramFoo,
	}

	m.MetadataDefs = append(m.MetadataDefs, diCompileUnit)
	m.MetadataDefs = append(m.MetadataDefs, diFile)
	m.MetadataDefs = append(m.MetadataDefs, emptyTuple)
	m.MetadataDefs = append(m.MetadataDefs, dwarfVersion)
	m.MetadataDefs = append(m.MetadataDefs, debugInfoVersion)
	m.MetadataDefs = append(m.MetadataDefs, wcharSize)
	m.MetadataDefs = append(m.MetadataDefs, picLevel)
	m.MetadataDefs = append(m.MetadataDefs, pieLevel)
	m.MetadataDefs = append(m.MetadataDefs, clangVersion)
	m.MetadataDefs = append(m.MetadataDefs, diSubprogramFoo)
	m.MetadataDefs = append(m.MetadataDefs, diSubroutineType)
	m.MetadataDefs = append(m.MetadataDefs, typesTuple)
	m.MetadataDefs = append(m.MetadataDefs, diBasicTypeI32)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarA)
	m.MetadataDefs = append(m.MetadataDefs, diLocA)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarB)
	m.MetadataDefs = append(m.MetadataDefs, diLocB)
	m.MetadataDefs = append(m.MetadataDefs, diLocalVarSum)
	m.MetadataDefs = append(m.MetadataDefs, diLocSum)

	// Named metadata definitions.
	//    !llvm.dbg.cu = !{!0}
	llvmDbgCu := &metadata.NamedDef{
		Name:  "llvm.dbg.cu",
		Nodes: []metadata.Node{diCompileUnit},
	}
	m.NamedMetadataDefs["llvm.dbg.cu"] = llvmDbgCu
	//    !llvm.module.flags = !{!3, !4, !5, !6, !7}
	llvmModuleFlags := &metadata.NamedDef{
		Name:  "llvm.module.flags",
		Nodes: []metadata.Node{dwarfVersion, debugInfoVersion, wcharSize, picLevel, pieLevel},
	}
	m.NamedMetadataDefs["llvm.module.flags"] = llvmModuleFlags
	//    !llvm.ident = !{!8}
	llvmIdent := &metadata.NamedDef{
		Name:  "llvm.ident",
		Nodes: []metadata.Node{clangVersion},
	}
	m.NamedMetadataDefs["llvm.ident"] = llvmIdent

	return diLocalVarA, diLocalVarB, diLocalVarSum, diLocA, diLocB, diLocSum
}

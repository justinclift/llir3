declare void @llvm.dbg.declare(metadata, metadata, metadata)

define i32 @foo(i32, i32) {
; <label>:2
	%3 = alloca i32
	%4 = alloca i32
	%5 = alloca i32
	store i32 %0, i32* %3
	call void @llvm.dbg.declare(metadata i32* %3, metadata !13, metadata !DIExpression()), !dbg !14
	store i32 %1, i32* %4
	call void @llvm.dbg.declare(metadata i32* %4, metadata !15, metadata !DIExpression()), !dbg !16
	call void @llvm.dbg.declare(metadata i32* %5, metadata !17, metadata !DIExpression()), !dbg !18
	%6 = load i32, i32* %3
	%7 = load i32, i32* %4
	%8 = add i32 %6, %7
	store i32 %8, i32* %5
	%9 = load i32, i32* %5
	ret i32 %9
}

define i32 @main() {
; <label>:0
	%1 = call i32 @foo(i32 12, i32 30)
	ret i32 %1
}

!llvm.dbg.cu = !{!0}
!llvm.ident = !{!8}
!llvm.module.flags = !{!3, !4, !5, !6, !7}

!0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "clang version 8.0.0 (tags/RELEASE_800/final)", emissionKind: FullDebug, enums: !2)
!1 = !DIFile(filename: "foo.c", directory: "/home/u/Desktop/foo")
!2 = !{}
!3 = !{i32 2, !"Dwarf Version", i32 4}
!4 = !{i32 2, !"Debug Info Version", i32 3}
!5 = !{i32 1, !"wchar_size", i32 4}
!6 = !{i32 7, !"PIC Level", i32 2}
!7 = !{i32 7, !"PIE Level", i32 2}
!8 = !{!"clang version 8.0.0 (tags/RELEASE_800/final)"}
!9 = distinct !DISubprogram(name: "foo", scope: !1, file: !1, line: 1, type: !10, scopeLine: 1, flags: DIFlagPrototyped, spFlags: DISPFlagDefinition, unit: !0, retainedNodes: !2)
!10 = !DISubroutineType(types: !11)
!11 = !{!12, !12, !12}
!12 = !DIBasicType(name: "int", size: 32, encoding: DW_ATE_signed)
!13 = !DILocalVariable(name: "a", arg: 1, scope: !9, file: !1, line: 1, type: !12)
!14 = !DILocation(line: 1, column: 13, scope: !9)
!15 = !DILocalVariable(name: "b", arg: 2, scope: !9, file: !1, line: 1, type: !12)
!16 = !DILocation(line: 1, column: 20, scope: !9)
!17 = !DILocalVariable(name: "sum", scope: !9, file: !1, line: 2, type: !12)
!18 = !DILocation(line: 2, column: 6, scope: !9)


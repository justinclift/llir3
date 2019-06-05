declare void @llvm.dbg.declare(metadata, metadata, metadata)

define i32 @foo(i32, i32) !dbg !7 {
; <label>:2
	%3 = alloca i32, align 4
	%4 = alloca i32, align 4
	%5 = alloca i32, align 4
	store i32 %0, i32* %3, align 4
	call void @llvm.dbg.declare(metadata i32* %3, metadata !11, metadata !DIExpression()), !dbg !12
	store i32 %1, i32* %4, align 4
	call void @llvm.dbg.declare(metadata i32* %4, metadata !13, metadata !DIExpression()), !dbg !14
	call void @llvm.dbg.declare(metadata i32* %5, metadata !15, metadata !DIExpression()), !dbg !16
	%6 = load i32, i32* %3, align 4, !dbg !17
	%7 = load i32, i32* %4, align 4, !dbg !18
	%8 = add nsw i32 %6, %7, !dbg !19
	store i32 %8, i32* %5, align 4, !dbg !20
	%9 = load i32, i32* %5, align 4, !dbg !21
	ret i32 %9, !dbg !22
}

define i32 @main() !dbg !23 {
; <label>:0
	%1 = alloca i32, align 4
	store i32 0, i32* %1, align 4
	%2 = call i32 @foo(i32 12, i32 30), !dbg !26
	ret i32 %2, !dbg !27
}

!llvm.dbg.cu = !{!0}
!llvm.ident = !{!6}
!llvm.module.flags = !{!3, !4, !5}

!0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git 295e7b4486abfc0f01e6d22276931165f324c61e)", emissionKind: FullDebug, enums: !2, nameTableKind: None)
!1 = !DIFile(filename: "target.c", directory: "/home/jc/git_repos2/llir3")
!2 = !{}
!3 = !{i32 2, !"Dwarf Version", i32 4}
!4 = !{i32 2, !"Debug Info Version", i32 3}
!5 = !{i32 1, !"wchar_size", i32 4}
!6 = !{!"clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git 295e7b4486abfc0f01e6d22276931165f324c61e)"}
!7 = distinct !DISubprogram(name: "foo", scope: !1, file: !1, line: 1, type: !8, scopeLine: 1, flags: DIFlagPrototyped, spFlags: DISPFlagDefinition, unit: !0, retainedNodes: !2)
!8 = !DISubroutineType(types: !9)
!9 = !{!10, !10, !10}
!10 = !DIBasicType(name: "int", size: 32, encoding: DW_ATE_signed)
!11 = !DILocalVariable(name: "a", arg: 1, scope: !7, file: !1, line: 1, type: !10)
!12 = !DILocation(line: 1, column: 13, scope: !7)
!13 = !DILocalVariable(name: "b", arg: 2, scope: !7, file: !1, line: 1, type: !10)
!14 = !DILocation(line: 1, column: 20, scope: !7)
!15 = !DILocalVariable(name: "sum", scope: !7, file: !1, line: 2, type: !10)
!16 = !DILocation(line: 2, column: 9, scope: !7)
!17 = !DILocation(line: 3, column: 11, scope: !7)
!18 = !DILocation(line: 3, column: 15, scope: !7)
!19 = !DILocation(line: 3, column: 13, scope: !7)
!20 = !DILocation(line: 3, column: 9, scope: !7)
!21 = !DILocation(line: 4, column: 12, scope: !7)
!22 = !DILocation(line: 4, column: 5, scope: !7)
!23 = distinct !DISubprogram(name: "main", scope: !1, file: !1, line: 7, type: !24, scopeLine: 7, spFlags: DISPFlagDefinition, unit: !0, retainedNodes: !2)
!24 = !DISubroutineType(types: !25)
!25 = !{!10}
!26 = !DILocation(line: 8, column: 12, scope: !23)
!27 = !DILocation(line: 8, column: 5, scope: !23)


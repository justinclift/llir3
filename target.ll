; ModuleID = 'target.c'
source_filename = "target.c"
target datalayout = "e-m:e-i64:64-f80:128-n8:16:32:64-S128"
target triple = "x86_64-pc-linux-gnu"

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @foo(i32, i32) #0 !dbg !7 {
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

; Function Attrs: nounwind readnone speculatable
declare void @llvm.dbg.declare(metadata, metadata, metadata) #1

; Function Attrs: noinline nounwind optnone uwtable
define dso_local i32 @main() #0 !dbg !23 {
  %1 = alloca i32, align 4
  store i32 0, i32* %1, align 4
  %2 = call i32 @foo(i32 12, i32 30), !dbg !26
  ret i32 %2, !dbg !27
}

attributes #0 = { noinline nounwind optnone uwtable "correctly-rounded-divide-sqrt-fp-math"="false" "disable-tail-calls"="false" "less-precise-fpmad"="false" "min-legal-vector-width"="0" "no-frame-pointer-elim"="true" "no-frame-pointer-elim-non-leaf" "no-infs-fp-math"="false" "no-jump-tables"="false" "no-nans-fp-math"="false" "no-signed-zeros-fp-math"="false" "no-trapping-math"="false" "stack-protector-buffer-size"="8" "target-cpu"="x86-64" "target-features"="+fxsr,+mmx,+sse,+sse2,+x87" "unsafe-fp-math"="false" "use-soft-float"="false" }
attributes #1 = { nounwind readnone speculatable }

!llvm.dbg.cu = !{!0}
!llvm.module.flags = !{!3, !4, !5}
!llvm.ident = !{!6}

!0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "clang version 8.0.1 (https://git.llvm.org/git/clang.git ccfe04576c13497b9c422ceef0b6efe99077a392) (https://git.llvm.org/git/llvm.git 295e7b4486abfc0f01e6d22276931165f324c61e)", isOptimized: false, runtimeVersion: 0, emissionKind: FullDebug, enums: !2, nameTableKind: None)
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

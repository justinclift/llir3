; ModuleID = 'main.go'
source_filename = "main.go"
target datalayout = "e-m:e-p:32:32-i64:64-n32:64-S128"
target triple = "wasm32-unknown-unknown-wasm"

@runtime.runqueueBack = internal unnamed_addr global i8* null
@runtime.runqueueFront = internal unnamed_addr global i8* null
@runtime.stdout = internal unnamed_addr global i32 0
@runtime.nilPanic$string = internal unnamed_addr constant [23 x i8] c"nil pointer dereference"
@runtime.runtimePanic$string = internal unnamed_addr constant [22 x i8] c"panic: runtime error: "
@main.go.main$string = internal unnamed_addr constant [12 x i8] c"Hello world!"

; Function Attrs: optsize
define void @_start() local_unnamed_addr #0 section ".text._start" !dbg !5 {
entry:
  %0 = tail call i32 @io_get_stdout(), !dbg !8
  store i32 %0, i32* @runtime.stdout, align 4, !dbg !8
  ret void, !dbg !12
}

; Function Attrs: optsize
define dso_local void @runtime.activateTask(i8*, i8* nocapture readnone %context, i8* nocapture readnone %parentHandle) unnamed_addr #0 section ".text.runtime.activateTask" !dbg !13 {
entry:
  call void @llvm.dbg.value(metadata i8* %0, metadata !21, metadata !DIExpression()), !dbg !22
  %1 = icmp eq i8* %0, null, !dbg !23
  br i1 %1, label %if.then, label %if.done, !dbg !24

if.then:                                          ; preds = %if.then.i, %store.next.i, %if.then4.i, %entry
  ret void, !dbg !25

if.done:                                          ; preds = %entry
  call void @llvm.dbg.value(metadata i8* %0, metadata !26, metadata !DIExpression()), !dbg !29
  %2 = bitcast i8* %0 to i8**, !dbg !31
  %3 = load i8*, i8** %2, align 4, !dbg !31
  %4 = icmp eq i8* %3, null, !dbg !31
  br i1 %4, label %if.then.i, label %if.done3.i, !dbg !32

if.then.i:                                        ; preds = %if.done
  %5 = bitcast i8* %0 to { i8*, i8* }*
  %6 = getelementptr inbounds { i8*, i8* }, { i8*, i8* }* %5, i32 0, i32 1
  %7 = load i8*, i8** %6
  %8 = bitcast i8* %7 to void (i8*)*
  tail call fastcc void %8(i8* nonnull %0), !dbg !33
  br label %if.then, !dbg !34

if.done3.i:                                       ; preds = %if.done
  %9 = load i8*, i8** @runtime.runqueueBack, align 4, !dbg !35
  %10 = icmp eq i8* %9, null, !dbg !36
  br i1 %10, label %if.then4.i, label %store.next.i, !dbg !32

if.then4.i:                                       ; preds = %if.done3.i
  store i8* %0, i8** @runtime.runqueueBack, align 4, !dbg !37
  store i8* %0, i8** @runtime.runqueueFront, align 4, !dbg !38
  br label %if.then, !dbg !32

store.next.i:                                     ; preds = %if.done3.i
  call void @llvm.dbg.value(metadata i8* %9, metadata !39, metadata !DIExpression()), !dbg !42
  %11 = getelementptr inbounds i8, i8* %9, i32 8, !dbg !44
  %12 = bitcast i8* %11 to i8**, !dbg !45
  store i8* %0, i8** %12, align 4, !dbg !45
  store i8* %0, i8** @runtime.runqueueBack, align 4, !dbg !46
  br label %if.then, !dbg !32
}

; Function Attrs: optsize
define void @cwa_main() local_unnamed_addr #0 section ".text.cwa_main" !dbg !47 {
entry:
  %0 = tail call i32 @io_get_stdout(), !dbg !48
  store i32 %0, i32* @runtime.stdout, align 4, !dbg !48
  tail call fastcc void @runtime.printstring(i8* getelementptr inbounds ([12 x i8], [12 x i8]* @main.go.main$string, i32 0, i32 0), i32 12), !dbg !50
  tail call fastcc void @runtime.printnl(), !dbg !50
  ret void, !dbg !53
}

; Function Attrs: noreturn optsize
define internal fastcc void @runtime.getFuncPtr() unnamed_addr #1 section ".text.runtime.getFuncPtr" !dbg !54 {
entry:
  call void @llvm.dbg.value(metadata i8* null, metadata !67, metadata !DIExpression(DW_OP_LLVM_fragment, 0, 32)), !dbg !69
  call void @llvm.dbg.value(metadata i32 0, metadata !67, metadata !DIExpression(DW_OP_LLVM_fragment, 32, 32)), !dbg !69
  call void @llvm.dbg.value(metadata i8* undef, metadata !68, metadata !DIExpression()), !dbg !69
  tail call fastcc void @runtime.nilPanic(), !dbg !70
  unreachable, !dbg !70
}

; Function Attrs: optsize
declare i32 @io_get_stdout() local_unnamed_addr #0

; Function Attrs: optsize
define i8* @memset(i8* nocapture returned, i8, i32) local_unnamed_addr #0 section ".text.memset" !dbg !71 {
entry:
  call void @llvm.dbg.value(metadata i8* %0, metadata !76, metadata !DIExpression()), !dbg !79
  call void @llvm.dbg.value(metadata i8 %1, metadata !77, metadata !DIExpression()), !dbg !79
  call void @llvm.dbg.value(metadata i32 %2, metadata !78, metadata !DIExpression()), !dbg !79
  br label %for.loop, !dbg !80

for.loop:                                         ; preds = %store.next, %entry
  %3 = phi i32 [ 0, %entry ], [ %7, %store.next ], !dbg !81
  %4 = icmp ult i32 %3, %2, !dbg !82
  br i1 %4, label %for.body, label %for.done, !dbg !80

for.body:                                         ; preds = %for.loop
  %5 = getelementptr inbounds i8, i8* %0, i32 %3, !dbg !83
  %6 = icmp eq i8* %5, null, !dbg !84
  br i1 %6, label %store.nil, label %store.next, !dbg !84

for.done:                                         ; preds = %for.loop
  ret i8* %0, !dbg !85

store.nil:                                        ; preds = %for.body
  tail call fastcc void @runtime.nilPanic(), !dbg !84
  unreachable, !dbg !84

store.next:                                       ; preds = %for.body
  store i8 %1, i8* %5, align 1, !dbg !84
  %7 = add i32 %3, 1, !dbg !86
  br label %for.loop, !dbg !80
}

; Function Attrs: noreturn optsize
define internal fastcc void @runtime.nilPanic() unnamed_addr #1 section ".text.runtime.nilPanic" !dbg !87 {
entry:
  tail call fastcc void @runtime.runtimePanic(), !dbg !89
  unreachable
}

; Function Attrs: optsize
define internal fastcc void @runtime.printnl() unnamed_addr #0 section ".text.runtime.printnl" !dbg !90 {
entry:
  tail call fastcc void @runtime.putchar(i8 13), !dbg !92
  tail call fastcc void @runtime.putchar(i8 10), !dbg !93
  ret void, !dbg !94
}

; Function Attrs: optsize
define internal fastcc void @runtime.printstring(i8* nocapture readonly, i32) unnamed_addr #0 section ".text.runtime.printstring" !dbg !95 {
entry:
  call void @llvm.dbg.value(metadata i8* %0, metadata !103, metadata !DIExpression(DW_OP_LLVM_fragment, 0, 32)), !dbg !104
  call void @llvm.dbg.value(metadata i32 %1, metadata !103, metadata !DIExpression(DW_OP_LLVM_fragment, 32, 32)), !dbg !104
  br label %for.loop, !dbg !105

for.loop:                                         ; preds = %for.body, %entry
  %2 = phi i32 [ 0, %entry ], [ %6, %for.body ], !dbg !106
  %3 = icmp slt i32 %2, %1, !dbg !107
  br i1 %3, label %for.body, label %for.done, !dbg !105

for.body:                                         ; preds = %for.loop
  %4 = getelementptr inbounds i8, i8* %0, i32 %2, !dbg !108
  %5 = load i8, i8* %4, align 1, !dbg !108
  tail call fastcc void @runtime.putchar(i8 %5), !dbg !109
  %6 = add nuw i32 %2, 1, !dbg !110
  br label %for.loop, !dbg !105

for.done:                                         ; preds = %for.loop
  ret void, !dbg !105
}

; Function Attrs: optsize
define internal fastcc void @runtime.putchar(i8) unnamed_addr #0 section ".text.runtime.putchar" !dbg !111 {
entry:
  %stackalloc.alloca = alloca [1 x i32], align 4, !dbg !116
  %.fca.0.gep = getelementptr inbounds [1 x i32], [1 x i32]* %stackalloc.alloca, i32 0, i32 0, !dbg !116
  store i32 0, i32* %.fca.0.gep, align 4, !dbg !116
  %stackalloc = bitcast [1 x i32]* %stackalloc.alloca to i8*, !dbg !116
  call void @llvm.dbg.value(metadata i8 %0, metadata !115, metadata !DIExpression()), !dbg !116
  store i8 %0, i8* %stackalloc, align 4, !dbg !117
  %1 = load i32, i32* @runtime.stdout, align 4, !dbg !118
  %2 = call i32 @resource_write(i32 %1, i8* nonnull %stackalloc, i32 1), !dbg !119
  ret void, !dbg !117
}

; Function Attrs: optsize
declare i32 @resource_write(i32, i8* nocapture, i32) local_unnamed_addr #0

; Function Attrs: noreturn optsize
define void @resume() local_unnamed_addr #1 section ".text.resume" !dbg !120 {
entry:
  tail call fastcc void @runtime.getFuncPtr(), !dbg !121
  unreachable
}

; Function Attrs: noreturn optsize
define internal fastcc void @runtime.runtimePanic() unnamed_addr #1 section ".text.runtime.runtimePanic" !dbg !122 {
entry:
  call void @llvm.dbg.value(metadata i8* getelementptr inbounds ([23 x i8], [23 x i8]* @runtime.nilPanic$string, i32 0, i32 0), metadata !124, metadata !DIExpression(DW_OP_LLVM_fragment, 0, 32)), !dbg !125
  call void @llvm.dbg.value(metadata i32 23, metadata !124, metadata !DIExpression(DW_OP_LLVM_fragment, 32, 32)), !dbg !125
  tail call fastcc void @runtime.printstring(i8* getelementptr inbounds ([22 x i8], [22 x i8]* @runtime.runtimePanic$string, i32 0, i32 0), i32 22), !dbg !126
  tail call fastcc void @runtime.printstring(i8* getelementptr inbounds ([23 x i8], [23 x i8]* @runtime.nilPanic$string, i32 0, i32 0), i32 23), !dbg !127
  tail call fastcc void @runtime.printnl(), !dbg !127
  tail call void @llvm.trap() #5, !dbg !128
  unreachable, !dbg !131
}

; Function Attrs: cold noreturn nounwind optsize
declare void @llvm.trap() #2

; Function Attrs: nounwind optsize readnone speculatable
declare void @llvm.dbg.value(metadata, metadata, metadata) #3

; Function Attrs: argmemonly nounwind optsize readonly
declare i8* @llvm.coro.subfn.addr(i8* nocapture readonly, i8) #4

attributes #0 = { optsize }
attributes #1 = { noreturn optsize }
attributes #2 = { cold noreturn nounwind optsize }
attributes #3 = { nounwind optsize readnone speculatable }
attributes #4 = { argmemonly nounwind optsize readonly }
attributes #5 = { nounwind }

!llvm.dbg.cu = !{!0}
!llvm.module.flags = !{!3, !4}

!0 = distinct !DICompileUnit(language: DW_LANG_C99, file: !1, producer: "TinyGo", isOptimized: true, runtimeVersion: 0, emissionKind: FullDebug, enums: !2)
!1 = !DIFile(filename: "main.go", directory: "")
!2 = !{}
!3 = !{i32 1, !"Debug Info Version", i32 3}
!4 = !{i32 1, !"Dwarf Version", i32 4}
!5 = distinct !DISubprogram(name: "runtime._start", linkageName: "_start", scope: !6, file: !6, line: 26, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!6 = !DIFile(filename: "runtime_wasm.go", directory: "../../../runtime")
!7 = !DISubroutineType(types: !2)
!8 = !DILocation(line: 11, column: 6, scope: !9, inlinedAt: !11)
!9 = distinct !DISubprogram(name: "runtime.initAll", linkageName: "runtime.initAll", scope: !10, file: !10, line: 11, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!10 = !DIFile(filename: "runtime.go", directory: "../../../runtime")
!11 = distinct !DILocation(line: 27, column: 9, scope: !5)
!12 = !DILocation(line: 0, scope: !5)
!13 = distinct !DISubprogram(name: "runtime.activateTask", linkageName: "runtime.activateTask", scope: !14, file: !14, line: 106, type: !15, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !20)
!14 = !DIFile(filename: "scheduler.go", directory: "../../../runtime")
!15 = !DISubroutineType(types: !16)
!16 = !{!17}
!17 = !DIDerivedType(tag: DW_TAG_pointer_type, baseType: !18, size: 32, align: 32, dwarfAddressSpace: 0)
!18 = !DIDerivedType(tag: DW_TAG_typedef, name: "runtime.coroutine", baseType: !19)
!19 = !DIBasicType(name: "uint8", size: 8, encoding: DW_ATE_unsigned)
!20 = !{!21}
!21 = !DILocalVariable(name: "task", arg: 1, scope: !13, file: !14, line: 106, type: !17)
!22 = !DILocation(line: 106, column: 6, scope: !13)
!23 = !DILocation(line: 107, column: 10, scope: !13)
!24 = !DILocation(line: 0, scope: !13)
!25 = !DILocation(line: 108, column: 3, scope: !13)
!26 = !DILocalVariable(name: "t", arg: 1, scope: !27, file: !14, line: 137, type: !17)
!27 = distinct !DISubprogram(name: "runtime.runqueuePushBack", linkageName: "runtime.runqueuePushBack", scope: !14, file: !14, line: 137, type: !15, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !28)
!28 = !{!26}
!29 = !DILocation(line: 137, column: 6, scope: !27, inlinedAt: !30)
!30 = distinct !DILocation(line: 111, column: 18, scope: !13)
!31 = !DILocation(line: 138, column: 11, scope: !27, inlinedAt: !30)
!32 = !DILocation(line: 0, scope: !27, inlinedAt: !30)
!33 = !DILocation(line: 140, column: 12, scope: !27, inlinedAt: !30)
!34 = !DILocation(line: 141, column: 3, scope: !27, inlinedAt: !30)
!35 = !DILocation(line: 148, column: 5, scope: !27, inlinedAt: !30)
!36 = !DILocation(line: 148, column: 18, scope: !27, inlinedAt: !30)
!37 = !DILocation(line: 150, column: 3, scope: !27, inlinedAt: !30)
!38 = !DILocation(line: 151, column: 3, scope: !27, inlinedAt: !30)
!39 = !DILocalVariable(name: "t", arg: 1, scope: !40, file: !14, line: 46, type: !17)
!40 = distinct !DISubprogram(name: "(*runtime.coroutine).promise", linkageName: "(*runtime.coroutine).promise", scope: !14, file: !14, line: 46, type: !15, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !41)
!41 = !{!39}
!42 = !DILocation(line: 46, column: 21, scope: !40, inlinedAt: !43)
!43 = distinct !DILocation(line: 154, column: 42, scope: !27, inlinedAt: !30)
!44 = !DILocation(line: 47, column: 32, scope: !40, inlinedAt: !43)
!45 = !DILocation(line: 155, column: 19, scope: !27, inlinedAt: !30)
!46 = !DILocation(line: 156, column: 3, scope: !27, inlinedAt: !30)
!47 = distinct !DISubprogram(name: "runtime.cwa_main", linkageName: "cwa_main", scope: !6, file: !6, line: 31, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!48 = !DILocation(line: 11, column: 6, scope: !9, inlinedAt: !49)
!49 = distinct !DILocation(line: 32, column: 9, scope: !47)
!50 = !DILocation(line: 4, column: 9, scope: !51, inlinedAt: !52)
!51 = distinct !DISubprogram(name: "main.go.main", linkageName: "main.go.main", scope: !1, file: !1, line: 3, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!52 = distinct !DILocation(line: 33, column: 10, scope: !47)
!53 = !DILocation(line: 0, scope: !47)
!54 = distinct !DISubprogram(name: "runtime.getFuncPtr", linkageName: "runtime.getFuncPtr", scope: !55, file: !55, line: 26, type: !56, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !66)
!55 = !DIFile(filename: "func.go", directory: "../../../runtime")
!56 = !DISubroutineType(types: !57)
!57 = !{!58, !65}
!58 = !DIDerivedType(tag: DW_TAG_typedef, name: "runtime.funcValue", baseType: !59)
!59 = !DICompositeType(tag: DW_TAG_structure_type, size: 64, align: 32, elements: !60)
!60 = !{!61, !63}
!61 = !DIDerivedType(tag: DW_TAG_member, name: "context", baseType: !62, size: 32, align: 32)
!62 = !DIDerivedType(tag: DW_TAG_pointer_type, name: "unsafe.Pointer", baseType: null, size: 32, align: 32, dwarfAddressSpace: 0)
!63 = !DIDerivedType(tag: DW_TAG_member, name: "id", baseType: !64, size: 32, align: 32, offset: 32)
!64 = !DIBasicType(name: "uintptr", size: 32, encoding: DW_ATE_unsigned)
!65 = !DIDerivedType(tag: DW_TAG_pointer_type, baseType: !19, size: 32, align: 32, dwarfAddressSpace: 0)
!66 = !{!67, !68}
!67 = !DILocalVariable(name: "val", arg: 1, scope: !54, file: !55, line: 26, type: !58)
!68 = !DILocalVariable(name: "signature", arg: 2, scope: !54, file: !55, line: 26, type: !65)
!69 = !DILocation(line: 26, column: 6, scope: !54)
!70 = !DILocation(line: 27, column: 59, scope: !54)
!71 = distinct !DISubprogram(name: "runtime.memset", linkageName: "memset", scope: !6, file: !6, line: 73, type: !72, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !75)
!72 = !DISubroutineType(types: !73)
!73 = !{!62, !74, !64}
!74 = !DIBasicType(name: "byte", size: 8, encoding: DW_ATE_unsigned)
!75 = !{!76, !77, !78}
!76 = !DILocalVariable(name: "ptr", arg: 1, scope: !71, file: !6, line: 73, type: !62)
!77 = !DILocalVariable(name: "c", arg: 2, scope: !71, file: !6, line: 73, type: !74)
!78 = !DILocalVariable(name: "size", arg: 3, scope: !71, file: !6, line: 73, type: !64)
!79 = !DILocation(line: 73, column: 6, scope: !71)
!80 = !DILocation(line: 0, scope: !71)
!81 = !DILocation(line: 74, column: 6, scope: !71)
!82 = !DILocation(line: 74, column: 25, scope: !71)
!83 = !DILocation(line: 75, column: 26, scope: !71)
!84 = !DILocation(line: 75, column: 3, scope: !71)
!85 = !DILocation(line: 77, column: 2, scope: !71)
!86 = !DILocation(line: 74, column: 33, scope: !71)
!87 = distinct !DISubprogram(name: "runtime.nilPanic", linkageName: "runtime.nilPanic", scope: !88, file: !88, line: 39, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!88 = !DIFile(filename: "panic.go", directory: "../../../runtime")
!89 = !DILocation(line: 40, column: 14, scope: !87)
!90 = distinct !DISubprogram(name: "runtime.printnl", linkageName: "runtime.printnl", scope: !91, file: !91, line: 199, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!91 = !DIFile(filename: "print.go", directory: "../../../runtime")
!92 = !DILocation(line: 200, column: 9, scope: !90)
!93 = !DILocation(line: 201, column: 9, scope: !90)
!94 = !DILocation(line: 0, scope: !90)
!95 = distinct !DISubprogram(name: "runtime.printstring", linkageName: "runtime.printstring", scope: !91, file: !91, line: 12, type: !96, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !102)
!96 = !DISubroutineType(types: !97)
!97 = !{!98}
!98 = !DICompositeType(tag: DW_TAG_structure_type, name: "string", size: 64, align: 32, elements: !99)
!99 = !{!100, !101}
!100 = !DIDerivedType(tag: DW_TAG_member, name: "ptr", baseType: !65, size: 32, align: 32)
!101 = !DIDerivedType(tag: DW_TAG_member, name: "len", baseType: !64, size: 32, align: 32, offset: 32)
!102 = !{!103}
!103 = !DILocalVariable(name: "s", arg: 1, scope: !95, file: !91, line: 12, type: !98)
!104 = !DILocation(line: 12, column: 6, scope: !95)
!105 = !DILocation(line: 0, scope: !95)
!106 = !DILocation(line: 13, column: 6, scope: !95)
!107 = !DILocation(line: 13, column: 16, scope: !95)
!108 = !DILocation(line: 14, column: 12, scope: !95)
!109 = !DILocation(line: 14, column: 10, scope: !95)
!110 = !DILocation(line: 13, column: 26, scope: !95)
!111 = distinct !DISubprogram(name: "runtime.putchar", linkageName: "runtime.putchar", scope: !6, file: !6, line: 36, type: !112, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !114)
!112 = !DISubroutineType(types: !113)
!113 = !{!74}
!114 = !{!115}
!115 = !DILocalVariable(name: "c", arg: 1, scope: !111, file: !6, line: 36, type: !74)
!116 = !DILocation(line: 36, column: 6, scope: !111)
!117 = !DILocation(line: 0, scope: !111)
!118 = !DILocation(line: 37, column: 17, scope: !111)
!119 = !DILocation(line: 37, column: 16, scope: !111)
!120 = distinct !DISubprogram(name: "runtime.resume", linkageName: "resume", scope: !6, file: !6, line: 48, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!121 = !DILocation(line: 49, column: 13, scope: !120)
!122 = distinct !DISubprogram(name: "runtime.runtimePanic", linkageName: "runtime.runtimePanic", scope: !88, file: !88, line: 17, type: !96, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !123)
!123 = !{!124}
!124 = !DILocalVariable(name: "msg", arg: 1, scope: !122, file: !88, line: 17, type: !98)
!125 = !DILocation(line: 17, column: 6, scope: !122)
!126 = !DILocation(line: 18, column: 13, scope: !122)
!127 = !DILocation(line: 19, column: 9, scope: !122)
!128 = !DILocation(line: 69, column: 6, scope: !129, inlinedAt: !130)
!129 = distinct !DISubprogram(name: "runtime.abort", linkageName: "runtime.abort", scope: !6, file: !6, line: 68, type: !7, flags: DIFlagPrototyped, spFlags: DISPFlagLocalToUnit | DISPFlagDefinition | DISPFlagOptimized, unit: !0, retainedNodes: !2)
!130 = distinct !DILocation(line: 20, column: 7, scope: !122)
!131 = !DILocation(line: 0, scope: !129, inlinedAt: !130)

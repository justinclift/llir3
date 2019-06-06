(module
  (type $t0 (func))
  (type $t1 (func (result i32)))
  (type $t2 (func (param i32 i32 i32)))
  (type $t3 (func (param i32)))
  (type $t4 (func (param i32 i32)))
  (type $t5 (func (param i32 i32 i32) (result i32)))
  (import "env" "io_get_stdout" (func $env.io_get_stdout (type $t1)))
  (import "env" "resource_write" (func $env.resource_write (type $t5)))
  (func $f2 (type $t0)
    (local $l0 i32) (local $l1 i32)
    call $env.io_get_stdout
    local.set $l0
    i32.const 0
    local.set $l1
    local.get $l1
    local.get $l0
    i32.store offset=8
    return)
  (func $f3 (type $t2) (param $p0 i32) (param $p1 i32) (param $p2 i32)
    (local $l3 i32) (local $l4 i32) (local $l5 i32) (local $l6 i32) (local $l7 i32) (local $l8 i32) (local $l9 i32) (local $l10 i32) (local $l11 i32) (local $l12 i32) (local $l13 i32) (local $l14 i32) (local $l15 i32) (local $l16 i32) (local $l17 i32) (local $l18 i32) (local $l19 i32) (local $l20 i32) (local $l21 i32) (local $l22 i32) (local $l23 i32) (local $l24 i32) (local $l25 i32) (local $l26 i32) (local $l27 i32)
    i32.const 0
    local.set $l3
    local.get $p0
    local.set $l4
    local.get $l3
    local.set $l5
    local.get $l4
    local.get $l5
    i32.eq
    local.set $l6
    i32.const 1
    local.set $l7
    local.get $l6
    local.get $l7
    i32.and
    local.set $l8
    block $B0
      local.get $l8
      br_if $B0
      i32.const 0
      local.set $l9
      local.get $p0
      i32.load
      local.set $l10
      local.get $l10
      local.set $l11
      local.get $l9
      local.set $l12
      local.get $l11
      local.get $l12
      i32.eq
      local.set $l13
      i32.const 1
      local.set $l14
      local.get $l13
      local.get $l14
      i32.and
      local.set $l15
      block $B1
        local.get $l15
        i32.eqz
        br_if $B1
        local.get $p0
        i32.load offset=4
        local.set $l16
        local.get $p0
        local.get $l16
        call_indirect (type $t3) $T0
        br $B0
      end
      i32.const 0
      local.set $l17
      i32.const 0
      local.set $l18
      local.get $l18
      i32.load
      local.set $l19
      local.get $l19
      local.set $l20
      local.get $l17
      local.set $l21
      local.get $l20
      local.get $l21
      i32.eq
      local.set $l22
      i32.const 1
      local.set $l23
      local.get $l22
      local.get $l23
      i32.and
      local.set $l24
      block $B2
        local.get $l24
        i32.eqz
        br_if $B2
        i32.const 0
        local.set $l25
        local.get $l25
        local.get $p0
        i32.store
        i32.const 0
        local.set $l26
        local.get $l26
        local.get $p0
        i32.store offset=4
        br $B0
      end
      local.get $l19
      local.get $p0
      i32.store offset=8
      i32.const 0
      local.set $l27
      local.get $l27
      local.get $p0
      i32.store
    end
    return)
  (func $f4 (type $t0)
    (local $l0 i32) (local $l1 i32) (local $l2 i32) (local $l3 i32)
    i32.const 70
    local.set $l0
    i32.const 12
    local.set $l1
    call $env.io_get_stdout
    local.set $l2
    i32.const 0
    local.set $l3
    local.get $l3
    local.get $l2
    i32.store offset=8
    local.get $l0
    local.get $l1
    call $f5
    call $f6
    return)
  (func $f5 (type $t4) (param $p0 i32) (param $p1 i32)
    (local $l2 i32) (local $l3 i32) (local $l4 i32) (local $l5 i32) (local $l6 i32) (local $l7 i32) (local $l8 i32) (local $l9 i32) (local $l10 i32) (local $l11 i32) (local $l12 i32) (local $l13 i32)
    i32.const 0
    local.set $l2
    local.get $l2
    local.set $l3
    block $B0
      loop $L1
        local.get $l3
        local.set $l4
        local.get $l4
        local.set $l5
        local.get $p1
        local.set $l6
        local.get $l5
        local.get $l6
        i32.lt_s
        local.set $l7
        i32.const 1
        local.set $l8
        local.get $l7
        local.get $l8
        i32.and
        local.set $l9
        local.get $l9
        i32.eqz
        br_if $B0
        local.get $p0
        local.get $l4
        i32.add
        local.set $l10
        local.get $l10
        i32.load8_u
        local.set $l11
        local.get $l11
        call $f11
        i32.const 1
        local.set $l12
        local.get $l4
        local.get $l12
        i32.add
        local.set $l13
        local.get $l13
        local.set $l3
        br $L1
      end
    end
    return)
  (func $f6 (type $t0)
    (local $l0 i32) (local $l1 i32)
    i32.const 10
    local.set $l0
    i32.const 13
    local.set $l1
    local.get $l1
    call $f11
    local.get $l0
    call $f11
    return)
  (func $f7 (type $t0)
    call $f8
    unreachable)
  (func $f8 (type $t0)
    call $f10
    unreachable)
  (func $f9 (type $t5) (param $p0 i32) (param $p1 i32) (param $p2 i32) (result i32)
    (local $l3 i32) (local $l4 i32) (local $l5 i32) (local $l6 i32) (local $l7 i32) (local $l8 i32) (local $l9 i32) (local $l10 i32) (local $l11 i32) (local $l12 i32) (local $l13 i32) (local $l14 i32) (local $l15 i32) (local $l16 i32) (local $l17 i32) (local $l18 i32) (local $l19 i32)
    i32.const 0
    local.set $l3
    local.get $l3
    local.set $l4
    loop $L0 (result i32)
      local.get $l4
      local.set $l5
      local.get $l5
      local.set $l6
      local.get $p2
      local.set $l7
      local.get $l6
      local.get $l7
      i32.lt_u
      local.set $l8
      i32.const 1
      local.set $l9
      local.get $l8
      local.get $l9
      i32.and
      local.set $l10
      block $B1
        block $B2
          block $B3
            local.get $l10
            i32.eqz
            br_if $B3
            i32.const 0
            local.set $l11
            local.get $p0
            local.get $l5
            i32.add
            local.set $l12
            local.get $l12
            local.set $l13
            local.get $l11
            local.set $l14
            local.get $l13
            local.get $l14
            i32.eq
            local.set $l15
            i32.const 1
            local.set $l16
            local.get $l15
            local.get $l16
            i32.and
            local.set $l17
            local.get $l17
            br_if $B2
            br $B1
          end
          local.get $p0
          return
        end
        call $f8
        unreachable
      end
      local.get $l12
      local.get $p1
      i32.store8
      i32.const 1
      local.set $l18
      local.get $l5
      local.get $l18
      i32.add
      local.set $l19
      local.get $l19
      local.set $l4
      br $L0
    end)
  (func $f10 (type $t0)
    (local $l0 i32) (local $l1 i32) (local $l2 i32) (local $l3 i32)
    i32.const 16
    local.set $l0
    i32.const 23
    local.set $l1
    i32.const 48
    local.set $l2
    i32.const 22
    local.set $l3
    local.get $l2
    local.get $l3
    call $f5
    local.get $l0
    local.get $l1
    call $f5
    call $f6
    unreachable
    unreachable)
  (func $f11 (type $t3) (param $p0 i32)
    (local $l1 i32) (local $l2 i32) (local $l3 i32) (local $l4 i32) (local $l5 i32) (local $l6 i32) (local $l7 i32) (local $l8 i32) (local $l9 i32) (local $l10 i32) (local $l11 i32) (local $l12 i32)
    global.get $g0
    local.set $l1
    i32.const 16
    local.set $l2
    local.get $l1
    local.get $l2
    i32.sub
    local.set $l3
    local.get $l3
    global.set $g0
    i32.const 1
    local.set $l4
    i32.const 12
    local.set $l5
    local.get $l3
    local.get $l5
    i32.add
    local.set $l6
    local.get $l6
    local.set $l7
    i32.const 0
    local.set $l8
    local.get $l3
    local.get $l8
    i32.store offset=12
    local.get $l3
    local.get $p0
    i32.store8 offset=12
    i32.const 0
    local.set $l9
    local.get $l9
    i32.load offset=8
    local.set $l10
    local.get $l10
    local.get $l7
    local.get $l4
    call $env.resource_write
    drop
    i32.const 16
    local.set $l11
    local.get $l3
    local.get $l11
    i32.add
    local.set $l12
    local.get $l12
    global.set $g0
    return)
  (func $f12 (type $t0)
    call $f7
    unreachable)
  (table $T0 1 1 funcref)
  (memory $memory 2)
  (global $g0 (mut i32) (i32.const 66656))
  (export "memory" (memory 0))
  (export "main" (func $f4))
  (data $d0 (i32.const 0) "\00\00\00\00")
  (data $d1 (i32.const 4) "\00\00\00\00")
  (data $d2 (i32.const 8) "\00\00\00\00")
  (data $d3 (i32.const 16) "nil pointer dereference")
  (data $d4 (i32.const 48) "panic: runtime error: ")
  (data $d5 (i32.const 70) "Hello world!"))

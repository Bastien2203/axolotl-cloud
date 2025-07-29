import React from "react"
import { cn } from "../../libs/utils/cn"


type TabItem = {
  value: string
  label: string
  content: React.ReactNode
  disabled?: boolean
}

type TabsProps = {
  items: TabItem[]
  value?: string                    // contrôlé
  defaultValue?: string             // non contrôlé
  onValueChange?: (v: string) => void
  className?: string
}

export const Tabs: React.FC<TabsProps> = ({
  items,
  value,
  defaultValue,
  onValueChange,
  className
}) => {
  const isControlled = value !== undefined
  const [internal, setInternal] = React.useState<string>(defaultValue ?? items[0]?.value)
  const selected = isControlled ? value! : internal

  const setSelected = (v: string) => {
    if (!isControlled) setInternal(v)
    onValueChange?.(v)
  }

  const id = React.useId()
  const btnRefs = React.useRef<(HTMLButtonElement | null)[]>([])


  const indexOf = (v: string) => items.findIndex(i => i.value === v)

  const focusTab = (idx: number) => {
    const target = items[idx]
    if (!target || target.disabled) return
    setSelected(target.value)
    btnRefs.current[idx]?.focus()
  }

  const onKeyDown = (e: React.KeyboardEvent) => {
    const current = indexOf(selected)
    if (current < 0) return
    if (e.key === "ArrowRight") {
      e.preventDefault()
      for (let i = 1; i <= items.length; i++) {
        const next = (current + i) % items.length
        if (!items[next].disabled) return focusTab(next)
      }
    }
    if (e.key === "ArrowLeft") {
      e.preventDefault()
      for (let i = 1; i <= items.length; i++) {
        const prev = (current - i + items.length) % items.length
        if (!items[prev].disabled) return focusTab(prev)
      }
    }
    if (e.key === "Home") {
      e.preventDefault()
      const first = items.findIndex(i => !i.disabled)
      if (first >= 0) focusTab(first)
    }
    if (e.key === "End") {
      e.preventDefault()
      const last = [...items].reverse().findIndex(i => !i.disabled)
      if (last >= 0) focusTab(items.length - 1 - last)
    }
  }

  return (
    <div className={cn("w-full", className)}>
      {/* Tab list */}
      <div
        role="tablist"
        aria-orientation="horizontal"
        className="flex gap-2 border-b"
        onKeyDown={onKeyDown}
      >
        {items.map((t, i) => {
          const isActive = t.value === selected
          return (
            <button
              key={t.value}
              ref={el => { btnRefs.current[i] = el }}
              role="tab"
              id={`${id}-tab-${t.value}`}
              aria-controls={`${id}-panel-${t.value}`}
              aria-selected={isActive}
              disabled={t.disabled}
              onClick={() => !t.disabled && setSelected(t.value)}
              className={cn(
                "px-4 py-2 text-sm font-medium border-b-2 -mb-[1px] transition",
                "focus:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 rounded-t",
                t.disabled && "opacity-50 cursor-not-allowed",
                isActive
                  ? "border-blue-600 text-blue-700"
                  : "border-transparent text-gray-600 hover:text-gray-800 hover:border-gray-300"
              )}
            >
              {t.label}
            </button>
          )
        })}
      </div>

      {/* Panels */}
      {items.map(t => {
        const isActive = t.value === selected
        return (
          <div
            key={t.value}
            role="tabpanel"
            id={`${id}-panel-${t.value}`}
            aria-labelledby={`${id}-tab-${t.value}`}
            hidden={!isActive}
            className="pt-4"
          >
            {isActive && t.content}
          </div>
        )
      })}
    </div>
  )
}

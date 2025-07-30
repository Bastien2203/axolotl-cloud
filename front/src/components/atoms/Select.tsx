import React from "react"
import { cn } from "../../libs/utils/cn"

type SelectProps = {
  name: string
  value?: string
  onChange?: (e: React.ChangeEvent<HTMLSelectElement>) => void
  options: { label: string; value: string }[]
  required?: boolean
  disabled?: boolean
  className?: string
}

const Select: React.FC<SelectProps> = ({
  name,
  value,
  onChange,
  options,
  required = false,
  disabled = false,
  className = "",
}) => {
  return (
    <select
      id={name}
      name={name}
      value={value}
      onChange={onChange}
      required={required}
      disabled={disabled}
      className={cn(
        "px-4 py-2 border border-gray-300 rounded-xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition disabled:opacity-50 bg-white",
        className
      )}
    >
      {options.map((opt) => (
        <option key={opt.value} value={opt.value}>
          {opt.label}
        </option>
      ))}
    </select>
  )
}

export default Select

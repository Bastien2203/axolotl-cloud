import React from "react"
import { cn } from "../../libs/utils/cn"


type InputProps = {
  name: string
  type?: React.InputHTMLAttributes<HTMLInputElement>["type"]
  placeholder?: string
  value?: string
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void
  required?: boolean
  disabled?: boolean
  className?: string
}

const Input: React.FC<InputProps> = ({
  name,
  type = "text",
  placeholder = "",
  value,
  onChange,
  required = false,
  disabled = false,
  className = "",
}) => {
  return (
    <input
      id={name}
      name={name}
      type={type}
      placeholder={placeholder}
      value={value}
      onChange={onChange}
      required={required}
      disabled={disabled}
      className={cn(
        "px-4 py-2 border border-gray-300 rounded-xl text-sm outline-none focus:ring-2 focus:ring-blue-500 transition disabled:opacity-50",
        className
      )}
    />

  )
}

export default Input

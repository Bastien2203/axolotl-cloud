import React from "react"
import { cn } from "../../libs/utils/cn"
import { VariantsStyles, type Variant } from "../../libs/utils/variant"



type ButtonProps = {
  children: React.ReactNode
  type?: "button" | "submit" | "reset"
  onClick?: () => void
  disabled?: boolean
  variant?: Variant
  className?: string
}

const Button: React.FC<ButtonProps> = ({
  children,
  type = "button",
  onClick,
  disabled = false,
  variant = "primary",
  className = ""
}) => {
  const baseStyles =
    "px-4 py-2 rounded-xl text-sm font-semibold transition focus:outline-none disabled:opacity-50 disabled:pointer-events-none cursor-pointer flex items-center justify-center gap-2"

  

  return (
    <button
      type={type}
      onClick={onClick}
      disabled={disabled}
      className={cn(baseStyles, VariantsStyles[variant], className)}
    >
      {children}
    </button>
  )
}

export default Button

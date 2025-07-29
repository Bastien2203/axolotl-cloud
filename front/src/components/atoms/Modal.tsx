import React, { useEffect } from "react"
import { X } from "lucide-react"

const Modal = ({
  children,
  onClose
}: {
  children: React.ReactNode
  onClose: () => void
}) => {

  useEffect(() => {
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose()
    }
    document.addEventListener("keydown", handleKey)
    return () => document.removeEventListener("keydown", handleKey)
  }, [onClose])

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm"
      role="dialog"
      aria-modal="true"
      onClick={(e) => {
        if (e.target === e.currentTarget) onClose()
      }}
    >
      <div
        className="relative min-w-xl p-6 bg-white rounded-2xl shadow-xl animate-fade-in max-h-[90vh] overflow-y-auto"
        onClick={(e) => e.stopPropagation()}
      >
        <button
          className="absolute top-2 right-2 p-1 text-gray-600 hover:text-black cursor-pointer"
          onClick={onClose}
          aria-label="Close modal"
        >
          <X />
        </button>
        {children}
      </div>
    </div>
  )
}

export default Modal

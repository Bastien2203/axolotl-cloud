import React, { useState, useRef, useEffect, type ReactNode } from "react";
import { MoreVertical } from "lucide-react";
import type { Variant } from "../../libs/utils/variant";

export type Option = {
  label: string;
  onClick: () => void;
  variant?: Variant;
};

type Props = {
  children: ReactNode;
  options: Option[];
  absolute?: boolean; // Optional prop to center the menu
  className?: string;
  onClick?: () => void; 
};

const MoreMenu: React.FC<Props> = ({ children, options, absolute , className, onClick }) => {
  const [open, setOpen] = useState(false);
  const menuRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (!menuRef.current?.contains(event.target as Node)) {
        setOpen(false);
      }
    };

    if (open) document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [open]);

  return (
    <div className={`relative ${className}`} ref={menuRef}  onClick={onClick}>
        {children}
      <button
        aria-label="Open menu"
        onClick={(e) => {
          e.stopPropagation();
          e.preventDefault();
          setOpen((prev) => !prev);
        }}
        className={`z-20 p-2 rounded-md hover:bg-gray-100 transition-colors ${absolute ? "absolute top-2 right-2" : ""}`}
      >
        <MoreVertical className="w-5 h-5 text-gray-600" />
      </button>

     

      {open && (
        <div
          role="menu"
          className="absolute top-10 right-2 z-30 min-w-[8rem] bg-white border border-gray-200 rounded-md shadow-lg transition-all duration-150 ease-out animate-in fade-in zoom-in"
        >
          <ul className="py-1">
            {options.map(({ label, onClick, variant }, i) => (
              <li
                key={i}
                role="menuitem"
                onClick={(e) => {
                  e.stopPropagation();
                  e.preventDefault();
                  onClick();
                  setOpen(false);
                }}
                className={`px-4 py-2 text-sm cursor-pointer hover:bg-gray-100 transition-colors ${
                  variant === "danger" ? "text-red-600" : "text-gray-700"
                }`}
              >
                {label}
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
};

export default MoreMenu;

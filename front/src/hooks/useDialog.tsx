import { useState } from "react";

const Dialog = ({ isOpen, children }: { isOpen: boolean; children: React.ReactNode }) => {
  if (!isOpen) return null;
  return children;
};

export const useDialog = <T extends string>() => {
  const [openModals, setOpenModals] = useState<Record<T, boolean>>({} as Record<T, boolean>);

  const openDialog = (name: T) => setOpenModals((prev) => ({ ...prev, [name]: true }));
  const closeDialog = (name: T) => setOpenModals((prev) => ({ ...prev, [name]: false }));

  const dialog = (name: T, children: React.ReactNode) => (
    <Dialog isOpen={!!openModals[name]}>{children}</Dialog>
  );

  return { openDialog, closeDialog, dialog };
};

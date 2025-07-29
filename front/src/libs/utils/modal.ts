

export const hideModal = (id: string) => {
    const modal = document.getElementById(id) as HTMLDialogElement;
    if (!modal) return;
    modal.close();
};

export const showModal = (id: string) => {
    const modal = document.getElementById(id) as HTMLDialogElement;
    if (!modal) return;
    modal.showModal();
};
import Button, { type Variant } from "../atoms/Button";
import Modal from "../atoms/Modal";


const ValidationModal = ({onClose, variant, text, label, onConfirm}: {onClose: () => void, variant: Variant, text: string, label: string, onConfirm: () => void}) => {
    return <Modal onClose={onClose}>
        <div>
            <h2>Validation Required</h2>
            <p>{text}</p>
            <div className="flex justify-end gap-2 mt-4">
                <Button onClick={onClose} variant="secondary">Cancel</Button>
                <Button onClick={onConfirm} variant={variant}>{label}</Button>
            </div>
        </div>
    </Modal>
}

export default ValidationModal;
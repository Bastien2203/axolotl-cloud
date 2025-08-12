import { useState } from "react";
import Button from "../atoms/Button";
import Input from "../atoms/Input";
import Modal from "../atoms/Modal";
import { useToast } from "../../contexts/ToastContext";


type Field = {
    name: string;
    type: string;
    placeholder?: string;
    required?: boolean;
}

const FormModal = <T extends Record<string, any>>({ name, onClose, onSubmit, fields, defaultValue }: { name: string, onClose: () => void, onSubmit: (data: T) => void, fields: Field[], defaultValue?: T }) => {
    const [formData, setFormData] = useState<T>(defaultValue || {} as T);
    const toast = useToast();

    const validateData = (data: T) => {
        for (const field of fields) {
            if (field.required && !data[field.name]) {
                return false;
            }
        }
        return true;
    }

    return <Modal onClose={onClose}>
        <div className="min-w-[60vw]">
            <h2>{name}</h2>
            <br/>
            <form onSubmit={(e) => {
                e.preventDefault();
                if (validateData(formData)) {
                    onSubmit(formData)
                } else {
                    toast.error("Please fill in all required fields.");
                }
            }} className="flex flex-col gap-4">
                {
                    fields.map((field) => (
                        <Input
                            key={field.name}
                            type={field.type}
                            name={field.name}
                            className="w-full"
                            placeholder={field.placeholder ? `${field.placeholder} ${field.required ? "" : "(optional)"}` : ""}
                            required={field.required}
                            onChange={(e) => setFormData({ ...formData, [field.name]: e.target.value })}
                            value={(formData as any)[field.name] || ""}
                        />
                    ))
                }

                <hr className="my-4 border-gray-300" />

                <Button type="submit">
                    {defaultValue ? "Save" : "Submit"}
                </Button>
            </form>

        </div>
    </Modal>
}

export default FormModal;
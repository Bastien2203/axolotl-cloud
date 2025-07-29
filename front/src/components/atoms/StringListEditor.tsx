import { Trash, Plus } from "lucide-react";
import Button from "./Button";
import Input from "./Input";

type StringListEditorProps = {
  label: string;
  data: string[];
  onChange: (updated: string[]) => void;
  addLabel: string;
  variant?: "primary" | "secondary" | "danger";
  placeholder?: string;
};

const StringListEditor = ({ label, data, onChange, addLabel, variant, placeholder }: StringListEditorProps) => {
  return (
    <div>
      <h3>{label}</h3>
      <div className="flex flex-col gap-2">
        {data.map((value, i) => (
          <div className="flex gap-2 items-center" key={`${label}-${i}`}>
            <Input
              name={`${label}_${i}`}
              value={value}
              placeholder={placeholder}
              onChange={(e) => {
                const updated = [...data];
                updated[i] = e.target.value;
                onChange(updated);
              }}
              className="flex-1"
            />
            <Trash
              className="text-red-500 cursor-pointer flex-shrink-0"
              onClick={() => {
                const updated = data.filter((_, index) => index !== i);
                onChange(updated);
              }}
            />
          </div>
        ))}
      </div>
      <Button
        onClick={() => onChange([...data, ""])}
        variant={variant}
        className="my-2 gap-2 w-full"
      >
        <Plus />
        {addLabel}
      </Button>
    </div>
  );
};

export default StringListEditor;

import { useState } from "react";
import type { Project } from "../../api/types";
import Button from "../atoms/Button";
import Input from "../atoms/Input";
import Modal from "../atoms/Modal";

import { useToast } from "../../contexts/ToastContext";

type CreationProject = Omit<Project, "id" | "created_at" | "updated_at">;

const CreateProjectModal = ({ onClose, onCreate, defaultValue }: { onClose: () => void, onCreate: (project: CreationProject) => void, defaultValue?: CreationProject }) => {
    const [newProject, setNewProject] = useState<CreationProject>(defaultValue || {
        name: "",
        icon_url: "",
        website_url: "", // Optional field for project website URL
    });
    const toast = useToast();

    const validateForm = (project: CreationProject) => {
        if (!project.name || !project.icon_url) {
            return false;
        }
        return true;
    }

    return <Modal onClose={onClose}>
        <div className="min-w-[60vw]">
            <h2>Create New Project</h2>
            <form onSubmit={(e) => {
                e.preventDefault();
                if (validateForm(newProject)) {
                    onCreate(newProject)
                } else {
                    toast.error("Please fill in all required fields.");
                }
            }} className="flex flex-col gap-4">
                <Input type="text" name="name" className="w-full " placeholder="Project name" required onChange={(e) => setNewProject({ ...newProject, name: e.target.value })} value={newProject.name} />
                <Input type="text" name="icon_url" className="w-full" placeholder="Icon URL" onChange={(e) => setNewProject({ ...newProject, icon_url: e.target.value })} value={newProject.icon_url} />
                <Input type="url" name="website_url" className="w-full" placeholder="Website URL"  onChange={(e) => setNewProject({ ...newProject, website_url: e.target.value })} value={newProject.website_url} />

                <hr className="my-4 border-gray-300" />

                <Button type="submit">
                    {defaultValue ? "Save Changes" : "Create Project"}
                </Button>
            </form>

        </div>
    </Modal>
}

export default CreateProjectModal;
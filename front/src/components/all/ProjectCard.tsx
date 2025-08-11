import type { Project } from "../../api/types";
import { useDialog } from "../../hooks/useDialog";
import MoreMenu from "../atoms/MoreMenu";
import CreateProjectModal from "../modals/CreateProjectModal";
import ValidationModal from "../modals/ValidationModal";


const ProjectCard = ({ project, onClick, onDelete, onEdit }: { project: Project, onClick: () => void, onDelete: (project: Project) => void, onEdit: (project: Project) => void }) => {
    const {openDialog, closeDialog, dialog} = useDialog<"edit-project-modal" | "validation-modal-project-delete">();

    const handleDelete = () => {
        onDelete(project);
        closeDialog(`validation-modal-project-delete`);
    };

    const handleEdit = (updatedData: Omit<Project, 'id' | 'created_at' | 'updated_at'>) => {
        onEdit({ ...project, ...updatedData });
        closeDialog(`edit-project-modal`);
    };

    return (
        <>
            {dialog("validation-modal-project-delete", (
                <ValidationModal onClose={() => closeDialog("validation-modal-project-delete")} variant="danger" text="Are you sure you want to delete this project?" label="Delete" onConfirm={handleDelete} />
            ))}
            {dialog("edit-project-modal", (
                <CreateProjectModal onClose={() => closeDialog("edit-project-modal")} onCreate={handleEdit} defaultValue={project} />
            ))}
            <MoreMenu absolute options={project.website_url && project.website_url !== "" ? [
                GoToWebsiteOption(project),
                EditProjectOption(openDialog),
                DeleteProjectOption(openDialog)
            ] : [
                EditProjectOption(openDialog),
                DeleteProjectOption(openDialog)
            ]}>
                <div className="shadow p-4 rounded aspect-square flex flex-col items-center justify-between hover:opacity-80 hover:bg-gray-100 cursor-pointer relative bg-white" onClick={onClick}>
                    <div></div>
                    <img src={project.icon_url} alt={`${project.name} icon`} className="w-12 h-auto" />
                    <h3 >{project.name}</h3>
                </div>
            </MoreMenu>
        </>
    );
}


const GoToWebsiteOption = (project: Project) => ({
    label: "Open Website URL",
    onClick: () => {
        if (project.website_url) {
            window.open(project.website_url, "_blank");
        } else {
            console.warn("No website URL available for this project.");
        }
    }
})

const EditProjectOption = (openDialog: (id: any) => void) => ({
    label: "Edit Project",
    onClick: () => {
        openDialog(`edit-project-modal`);
    }
})

const DeleteProjectOption = (openDialog: (id: any) => void) => ({
    label: "Delete Project",
    onClick: () => {
        openDialog(`validation-modal-project-delete`);
    },
    variant: "danger"
})

export default ProjectCard;
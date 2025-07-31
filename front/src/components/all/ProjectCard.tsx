import type { Project } from "../../api/types";
import { hideModal, showModal } from "../../libs/utils/modal";
import MoreMenu from "../atoms/MoreMenu";
import CreateProjectModal from "../modals/CreateProjectModal";
import ValidationModal from "../modals/ValidationModal";


const ProjectCard = ({ project, onClick, onDelete, onEdit }: { project: Project, onClick: () => void, onDelete: (project: Project) => void, onEdit: (project: Project) => void }) => {

    const handleDelete = () => {
        onDelete(project);
        hideModal(`validation-modal-project-delete-${project.id}`);
    };

    const handleEdit = (updatedData: Omit<Project, 'id' | 'created_at' | 'updated_at'>) => {
        onEdit({ ...project, ...updatedData });
        hideModal(`edit-project-modal-${project.id}`);
    };

    return (
        <>
            <dialog id={`validation-modal-project-delete-${project.id}`} open={false}>
                <ValidationModal onClose={() => {
                    hideModal(`validation-modal-project-delete-${project.id}`);
                }} variant="danger" text="Are you sure you want to delete this project?" label="Delete" onConfirm={handleDelete} />
            </dialog>
            <dialog id={`edit-project-modal-${project.id}`} open={false}>
                <CreateProjectModal onClose={() => hideModal(`edit-project-modal-${project.id}`)} onCreate={handleEdit} defaultValue={project} />
            </dialog>
            <MoreMenu absolute options={project.website_url && project.website_url !== "" ? [
                GoToWebsiteOption(project),
                EditProjectOption(project),
                DeleteProjectOption(project)
            ] : [
                EditProjectOption(project),
                DeleteProjectOption(project)
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

const EditProjectOption = (project: Project) => ({
    label: "Edit Project",
    onClick: () => {
        showModal(`edit-project-modal-${project.id}`);
    }
})

const DeleteProjectOption = (project: Project) => ({
    label: "Delete Project",
    onClick: () => {
        showModal(`validation-modal-project-delete-${project.id}`);
    },
    variant: "danger"
})

export default ProjectCard;
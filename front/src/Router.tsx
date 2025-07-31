import { createBrowserRouter, RouterProvider } from "react-router-dom"
import Base from "./components/layouts/Base"
import Projects from "./components/pages/Projects"
import ProjectDetails from "./components/pages/ProjectDetails"
import Jobs from "./components/pages/Jobs"
import JobDetails from "./components/pages/JobDetails"


const Router= () => {
    const router = createBrowserRouter([
        {
            path: "/",
            element: <Base/>,
            children: [
                {
                    index: true,
                    element: <Projects />
                },
                {
                    path: "/projects/:projectId",
                    element: <ProjectDetails />
                },
                {
                    path: "/jobs",
                    element: <Jobs />
                },
                {
                    path: "/jobs/:jobId",
                    element: <JobDetails />
                }
            ]
        },  
    ])

    return <RouterProvider router={router} />;
}

export default Router;
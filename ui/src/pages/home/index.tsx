import { AppSidebar } from "@/components/Sidebar";
import {
    SidebarProvider,
    SidebarTrigger,
} from "@/components/ui/sidebar"

export function HomePage() {
    return (
        <SidebarProvider>
            <AppSidebar user={{}} />
            <main>
                <SidebarTrigger />
                <div>This is main</div>
            </main>
        </SidebarProvider>
    )
}
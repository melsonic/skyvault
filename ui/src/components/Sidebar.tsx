import * as React from "react"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "@/components/ui/sidebar"
import { NavUser } from "@/components/nav-user"
import { Avatar } from "./ui/avatar"
import { Cloud, ChevronsUpDown, Home, Folder, ClockFading, FileVideoIcon, Trash2, Image } from "lucide-react"

function ApplicationHeader() {
  return (
    <SidebarMenuButton
      size="lg"
      className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
    >
      <Avatar className="h-8 w-8 rounded-lg flex items-center justify-center">
        <Cloud />
      </Avatar>
      <div className="grid flex-1 text-left text-lg leading-tight">
        <span className="truncate font-medium">SkyVault</span>
      </div>
      <ChevronsUpDown className="ml-auto size-4" />
    </SidebarMenuButton>
  )
}

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar> & { user: any }) {

  const data = React.useMemo(() => {
    return {
      navMain: [
        {
          title: "Main",
          url: "#",
          items: [
            {
              title: "Home",
              icon: Home,
              url: "#",
            },
            {
              title: "My Drive",
              icon: Folder,
              url: "#",
            },
            {
              title: "Recent",
              icon: ClockFading,
              url: "#",
            },
            {
              title: "Videos",
              icon: FileVideoIcon,
              url: "#",
            },
            {
              title: "Photos",
              icon: Image,
              url: "#"
            },
            {
              title: "Trash",
              icon: Trash2,
              url: "#",
            },
          ],
        },
      ],
    }
  }, [])

  const [activeItem, setActiveItem] = React.useState(0);

  return (
    <Sidebar {...props}>
      <SidebarHeader>
        <ApplicationHeader />
      </SidebarHeader>
      <SidebarContent>
        {data.navMain.map((item) => (
          <SidebarGroup key={item.title}>
            <SidebarGroupLabel>{item.title}</SidebarGroupLabel>
            <SidebarGroupContent>
              <SidebarMenu>
                {item.items.map((item, index) => (
                  <SidebarMenuItem key={item.title} onClick={() => setActiveItem(index)}>
                    <SidebarMenuButton asChild isActive={index === activeItem}>
                      <a href={item.url} className="flex">
                        <item.icon />
                        <span className="pl-2 font-normal text-base">{item.title}</span>
                      </a>
                    </SidebarMenuButton>
                  </SidebarMenuItem>
                ))}
              </SidebarMenu>
            </SidebarGroupContent>
          </SidebarGroup>
        ))}
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={props.user} />
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  )
}

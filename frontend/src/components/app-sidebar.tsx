import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from "./animate-ui/components/radix/sidebar";
import { Link, LinkOptions, useMatchRoute } from "@tanstack/react-router";
import { ChevronsUpDown, BadgeCheck, Bell, LogOut } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage } from "./ui/avatar";
import useIsMobile from "@/hooks/use-is-mobile";
import { fetchUser } from "@/lib/fetch-user-server-fn";
import { useServerFn } from "@tanstack/react-start";
import { useQuery } from "@tanstack/react-query";
import {
  DropdownMenu,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuContent,
  DropdownMenuTrigger,
} from "./animate-ui/components/radix/dropdown-menu";
import FluentHome12Filled from "~icons/fluent/home-12-filled";
import FluentFood16Filled from "~icons/fluent/food-16-filled?height=32&width=32px";

type NavData = {
  linkOptions: LinkOptions;
  label: string;
  icon: React.FC<React.SVGProps<SVGSVGElement>>;
};

const DATA: NavData[] = [
  {
    linkOptions: {
      to: "/",
    },
    label: "Home",
    icon: FluentHome12Filled,
  },
  {
    linkOptions: {
      to: "/foods",
    },
    label: "Foods",
    icon: FluentFood16Filled,
  },
];

const AppSidebar = () => {
  const { isMobile } = useIsMobile();
  const getUser = useServerFn(fetchUser);

  const { data: user } = useQuery({
    queryKey: ["user"],
    queryFn: getUser,
  });

  return (
    <Sidebar>
      <SidebarHeader>
        {/* Team Switcher */}
        <SidebarMenu>
          <SidebarMenuItem>
            {/* <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            > */}
            {/* <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <activeTeam.logo className="size-4" />
              </div> */}
            <Link
              to="/"
              className="grid flex-1 text-center text-xl leading-tight"
            >
              <span className="truncate font-semibold">WEEATE</span>
              {/* <span className="truncate text-xs">{activeTeam.plan}</span> */}
            </Link>
            {/* <ChevronsUpDown className="ml-auto" /> */}
            {/* </SidebarMenuButton> */}
          </SidebarMenuItem>
        </SidebarMenu>
        {/* Team Switcher */}
      </SidebarHeader>

      <SidebarContent>
        <SidebarGroup>
          <SidebarMenu>
            {DATA.map((item, i) => (
              <SidebarMenuItem key={i}>
                <SidebarMenuButton asChild size="lg">
                  <Link
                    {...item.linkOptions}
                    className="rounded-lg px-3 py-2 text-xl font-medium text-muted-foreground transition-colors hover:text-foreground [&.active]:bg-primary/10 [&.active]:text-primary"
                  >
                    <item.icon className="mr-3 size-6!" />
                    <span>{item.label}</span>
                  </Link>
                </SidebarMenuButton>
              </SidebarMenuItem>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter>
        {/* Nav User */}
        <SidebarMenu>
          <SidebarMenuItem>
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <SidebarMenuButton
                  size="lg"
                  className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <Avatar className="h-8 w-8 rounded-lg">
                    <AvatarImage
                      src={user?.app_metadata.avatar_url}
                      alt={user?.app_metadata.display_name}
                    />
                    <AvatarFallback className="rounded-lg">
                      {user?.email?.substring(0, 2).toUpperCase()}
                    </AvatarFallback>
                  </Avatar>
                  <div className="grid flex-1 text-left text-sm leading-tight">
                    <span className="truncate font-semibold">
                      {user?.app_metadata.display_name || "Unnamed User"}
                    </span>
                    <span className="truncate text-xs">{user?.email}</span>
                  </div>
                  <ChevronsUpDown className="ml-auto size-4" />
                </SidebarMenuButton>
              </DropdownMenuTrigger>
              <DropdownMenuContent
                className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                side={isMobile ? "bottom" : "right"}
                align="end"
                sideOffset={4}
              >
                <DropdownMenuLabel className="p-0 font-normal">
                  <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <Avatar className="h-8 w-8 rounded-lg">
                      <AvatarImage
                        src={user?.app_metadata.avatar_url}
                        alt={user?.app_metadata.display_name}
                      />
                      <AvatarFallback className="rounded-lg">
                        {user?.email?.substring(0, 2).toUpperCase()}
                      </AvatarFallback>
                    </Avatar>
                    <div className="grid flex-1 text-left text-sm leading-tight">
                      <span className="truncate font-semibold">
                        {user?.app_metadata.display_name || "Unnamed User"}
                      </span>
                      <span className="truncate text-xs">{user?.email}</span>
                    </div>
                  </div>
                </DropdownMenuLabel>
                <DropdownMenuSeparator />
                <DropdownMenuGroup>
                  <DropdownMenuItem>
                    <BadgeCheck />
                    Account
                  </DropdownMenuItem>
                  <DropdownMenuItem>
                    <Bell />
                    Notifications
                  </DropdownMenuItem>
                </DropdownMenuGroup>
                <DropdownMenuSeparator />
                <DropdownMenuItem>
                  <LogOut />
                  Log out
                </DropdownMenuItem>
              </DropdownMenuContent>
            </DropdownMenu>
          </SidebarMenuItem>
        </SidebarMenu>
        {/* Nav User */}
      </SidebarFooter>
      <SidebarRail />
    </Sidebar>
  );
};

export default AppSidebar;

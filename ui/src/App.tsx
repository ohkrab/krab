import { AppShell, Code, ColorScheme, ColorSchemeProvider, Group, Image, MantineProvider, Navbar, Switch, Text, ThemeIcon, UnstyledButton, useMantineTheme } from '@mantine/core';
import { IconActivityHeartbeat, IconArticle, IconDatabase, IconDeviceFloppy, IconHierarchy2, IconLock, IconMoonStars, IconSun, IconUser } from '@tabler/icons-react';
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from 'react';
import { useState } from 'react';
import DatabaseList from './components/database_list';
import { BrowserRouter, Outlet, Route, Routes } from 'react-router-dom';
import TablespaceList from './components/tablespace_list';

const queryClient = new QueryClient();

interface MainLinkProps {
  icon: React.ReactNode;
  href: string;
  label: string;
}

function MainLink({ icon, label, href }: MainLinkProps) {
  return (
    <UnstyledButton
      component="a"
      href={href}
      sx={(theme) => ({
        display: 'block',
        width: '100%',
        padding: theme.spacing.sm,
        borderRadius: theme.radius.sm,
        color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.colors.gray[7],

        '&:hover': {
          backgroundColor:
            theme.colorScheme === 'dark' ? theme.colors.teal[8] : theme.colors.teal[2],
          color: theme.colorScheme === 'dark' ? theme.colors.dark[7] : theme.colors.gray[7],
        },
      })}
    >
      <Group>
        <ThemeIcon size="md" variant="transparent">
          {icon}
        </ThemeIcon>
        <Text size="md">
          {label}
        </Text>
      </Group>
    </UnstyledButton>
  );
}
const data = [
  { icon: <IconDatabase />, label: 'Databases', href: "/databases" },
  { icon: <IconLock />, label: 'Locks', href: "#" },
  { icon: <IconDeviceFloppy />, label: 'Tablespaces', href: "/tablespaces" },
  { icon: <IconUser />, label: 'Roles', href: "#" },
  { icon: <IconArticle />, label: 'WAL', href: "#" },
  { icon: <IconHierarchy2 />, label: 'Replication', href: "#" },
  { icon: <IconActivityHeartbeat />, label: 'Diagnostics', href: "#" },
];

const MainLinks = () => {
  const links = data.map((link) => <MainLink {...link} key={link.label} />);
  return <div>{links}</div>;
}

function App() {
  const [colorScheme, setColorScheme] = useState<ColorScheme>('light');
  const toggleColorScheme = (value?: ColorScheme) =>
    setColorScheme(value || (colorScheme === 'dark' ? 'light' : 'dark'));
  const theme = useMantineTheme();

  return (
    <BrowserRouter>
      <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
        <MantineProvider withGlobalStyles withNormalizeCSS theme={{ colorScheme }} >
          <QueryClientProvider client={queryClient}>
            <AppShell
              padding="xl"
              navbar={
                <Navbar width={{ base: 300 }} p="md">
                  <Navbar.Section grow>
                    <Group className="" position="center" pt="md" pb="2rem" px="md">
                      <Image
                        width="auto"
                        height={60}
                        src="https://ohkrab.dev/images/logo.svg"
                        alt="Oh, Krab!"
                      />
                    </Group>
                    <MainLinks />
                  </Navbar.Section>

                  <Navbar.Section>
                    <Group position="center" my={30}>
                      <Switch
                        checked={colorScheme === 'dark'}
                        onChange={() => toggleColorScheme()}
                        size="lg"
                        onLabel={<IconSun color={theme.white} size="1.25rem" stroke={1.5} />}
                        offLabel={<IconMoonStars color={theme.colors.gray[6]} size="1.25rem" stroke={1.5} />}
                      />
                      <Code sx={{ fontWeight: 700 }}>v0.8.0</Code>
                    </Group>
                  </Navbar.Section>
                </Navbar>
              }
            >
              <Outlet />
              <Routes>
                <Route path="/" element={<DatabaseList />} />
                <Route path="/databases" element={<DatabaseList />} />
                <Route path="/tablespaces" element={<TablespaceList />} />
              </Routes>
            </AppShell>
          </QueryClientProvider>
        </MantineProvider>
      </ColorSchemeProvider>
    </BrowserRouter>
  );
}

export default App;

import { AppShell, Box, Group, Header, Image, Navbar, ScrollArea, Text, ThemeIcon, UnstyledButton } from '@mantine/core';
import { IconActivityHeartbeat, IconArticle, IconCopy, IconDatabase, IconDeviceFloppy, IconHierarchy2, IconLock, IconUser } from '@tabler/icons-react';
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from 'react';

const queryClient = new QueryClient();

interface MainLinkProps {
  icon: React.ReactNode;
  color: string;
  label: string;
}

function MainLink({ icon, color, label }: MainLinkProps) {
  return (
    <UnstyledButton
      sx={(theme) => ({
        display: 'block',
        width: '100%',
        padding: theme.spacing.sm,
        borderRadius: theme.radius.sm,
        color: theme.colorScheme === 'dark' ? theme.colors.dark[0] : theme.black,

        '&:hover': {
          backgroundColor:
            theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[0],
        },
      })}
    >
      <Group>
        <ThemeIcon size="md" color={color} variant="dark">
          {icon}
        </ThemeIcon>
        <Text size="md">{label}</Text>
      </Group>
    </UnstyledButton>
  );
}
const data = [
  { icon: <IconDatabase size="1.2rem" />, color: 'green', label: 'Databases' },
  { icon: <IconLock size="1.2rem" />, color: 'blue', label: 'Locks' },
  { icon: <IconDeviceFloppy size="1.2rem" />, color: 'teal', label: 'Tablespaces' },
  { icon: <IconUser size="1.2rem" />, color: 'violet', label: 'Roles' },
  { icon: <IconArticle size="1.2rem" />, color: 'grape', label: 'WAL' },
  { icon: <IconHierarchy2 size="1.2rem" />, color: 'grape', label: 'Replication' },
  { icon: <IconActivityHeartbeat size="1.2rem" />, color: 'grape', label: 'Diagnostics' },
];

const MainLinks = () => {
  const links = data.map((link) => <MainLink {...link} key={link.label} />);
  return <div>{links}</div>;
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppShell
        padding="xl"
        navbar={
          <Navbar width={{ base: 300 }} p="md">
            <Navbar.Section>
              <Box pb="md">
                <MainLinks />
              </Box>
            </Navbar.Section>
          </Navbar>
        }
        header={
          <Header height={60} p="md">
            <Group sx={{ height: '100%' }} px="xs" position="apart">
              <Image
                width="auto"
                height={30}
                src="https://ohkrab.dev/images/favicon.svg"
                alt="Oh, Krab!"
              />
            </Group>
          </Header>
        }
      >
        Content
      </AppShell>

    </QueryClientProvider>
  );
}

export default App;

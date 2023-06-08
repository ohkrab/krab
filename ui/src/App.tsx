import { AppShell, Code, ColorScheme, ColorSchemeProvider, Divider, Group, Image, MantineProvider, Navbar, Switch, Text, ThemeIcon, UnstyledButton, useMantineTheme } from '@mantine/core';
import { IconActivityHeartbeat, IconArticle, IconDatabase, IconDeviceFloppy, IconHierarchy2, IconLock, IconMoonStars, IconSun, IconUser } from '@tabler/icons-react';
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from 'react';
import { useState } from 'react';

const queryClient = new QueryClient();

interface MainLinkProps {
  icon: React.ReactNode;
  color: string;
  label: string;
}

function MainLink({ icon, label }: MainLinkProps) {
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
            theme.colorScheme === 'dark' ? theme.colors.dark[6] : theme.colors.gray[2],
        },
      })}
    >
      <Group>
        <ThemeIcon size="md" variant="transparent">
          {icon}
        </ThemeIcon>
        <Text size="md">{label}</Text>
      </Group>
    </UnstyledButton>
  );
}
const data = [
  { icon: <IconDatabase size="1.2rem" />, label: 'Databases' },
  { icon: <IconLock size="1.2rem" />, label: 'Locks' },
  { icon: <IconDeviceFloppy size="1.2rem" />, label: 'Tablespaces' },
  { icon: <IconUser size="1.2rem" />, label: 'Roles' },
  { icon: <IconArticle size="1.2rem" />, label: 'WAL' },
  { icon: <IconHierarchy2 size="1.2rem" />, label: 'Replication' },
  { icon: <IconActivityHeartbeat size="1.2rem" />, label: 'Diagnostics' },
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
    <ColorSchemeProvider colorScheme={colorScheme} toggleColorScheme={toggleColorScheme}>
      <MantineProvider withGlobalStyles withNormalizeCSS theme={{ colorScheme }} >
        <QueryClientProvider client={queryClient}>
          <AppShell
            padding="xl"
            navbar={
              <Navbar width={{ base: 300 }} p="md">
                <Navbar.Section grow>
                  <Group className="" position="apart" py="xxs" px="xs">
                    <Image
                      width="auto"
                      height={30}
                      src="https://ohkrab.dev/images/favicon.svg"
                      alt="Oh, Krab!"
                    />
                    <Code sx={{ fontWeight: 700 }}>v0.8.0</Code>
                  </Group>
                  <Divider my="sm" />
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
                  </Group>
                </Navbar.Section>
              </Navbar>
            }
          >
            Content
          </AppShell>

        </QueryClientProvider>
      </MantineProvider>
    </ColorSchemeProvider>
  );
}

export default App;

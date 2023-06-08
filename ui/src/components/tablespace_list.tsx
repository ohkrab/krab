import { Alert, Badge, Button, Group, Skeleton, Table, ThemeIcon } from '@mantine/core';
import { IconAlertCircle, IconCircleCheck, IconDatabase, IconExternalLink, IconInfinity, IconToggleLeft, IconToggleRight } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';

function TablespaceList() {
  const { isLoading, error, data } = useQuery(['tablespaces'], () =>
    fetch('/api/tablespaces').then(res => res.json())
  )

  if (error) return (
    <Alert icon={<IconAlertCircle size="1rem" />} title="Something went wrong" color="red.8" variant="filled">
      Error occurred: {error.message}
    </Alert>
  )

  let rows;
  if (isLoading) {
    rows = [0, 1, 2, 3, 4, 5, 6, 7, 8, 9].map((i) => (
      <tr key={i}>
        <td><Skeleton height={20} /></td>
        <td><Skeleton height={20} /></td>
        <td><Skeleton height={20} /></td>
        <td><Skeleton height={20} /></td>
      </tr>
    ));
  } else {
    rows = data.data.map((d) => (
      <tr key={d.ID}>
        <td>
          <Group>
            {d.name}
            {d.isTemplate && <Badge variant="outline" size="xs" color="teal.8">template</Badge>}
          </Group>
        </td>
        <td>{d.size}</td>
        <td>{d.location}</td>
        <td>
          <Button component="a" href={`/roles/${d.ownerID}`} compact color="teal" variant="outline" leftIcon={<IconExternalLink size="0.9rem" />}>
            {d.ownerName}
          </Button>
        </td>
      </tr>
    ));
  }

  return (
    <>
      <h2>Tablespaces</h2>
      <Table verticalSpacing="xs">
        <thead>
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th>Location</th>
            <th>Owner</th>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </Table>
    </>
  );
}

export default TablespaceList;

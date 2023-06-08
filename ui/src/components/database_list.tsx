import { Alert, Anchor, Badge, Box, Button, Group, Skeleton, Table } from '@mantine/core';
import { IconAlertCircle, IconExternalLink } from '@tabler/icons-react';
import { useQuery } from '@tanstack/react-query';
import { Link } from 'react-router-dom';

function DatabaseList() {
  const { isLoading, error, data } = useQuery(['databases'], () =>
    fetch('/api/databases').then(res => res.json())
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
            <Anchor href={`/databases/${d.ID}`} color="teal">
              {d.name}
            </Anchor>
            {d.isTemplate && <Badge variant="outline" size="xs" color="teal.7">template</Badge>}
          </Group>
        </td>
        <td>{d.size}</td>
        <td>
          <Button component="a" href={`/tablespaces/${d.tablespaceID}`} compact color="teal" variant="outline" leftIcon={<IconExternalLink size="0.9rem" />}>
            {d.tablespaceName}
          </Button>
        </td>
        <td>{d.connectionLimit === -1 ? "Unlimited" : d.connectionLimit}</td>
        <td>
          <Button component="a" href={`/roles/${d.ownerID}`} compact color="teal" variant="outline" leftIcon={<IconExternalLink size="0.9rem" />}>
            {d.ownerName}
          </Button>
        </td>
        <td>{d.encoding}</td>
        <td>{d.collation}</td>
        <td>{d.characterType}</td>
      </tr>
    ));
  }

  return (
    <>
      <h2>Databases</h2>
      <Table verticalSpacing="xs">
        <thead>
          <tr>
            <th>Name</th>
            <th>Size</th>
            <th>Tablespace</th>
            <th>Connection Limit</th>
            <th>Owner</th>
            <th>Encoding</th>
            <th>Collation</th>
            <th>Character Type</th>
          </tr>
        </thead>
        <tbody>{rows}</tbody>
      </Table>
    </>
  );
}

export default DatabaseList;

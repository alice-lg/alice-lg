
import { useNeighbors }
 from 'app/components/neighbors/Provider';


const Neighbors = ({filter}) => {
  const {isLoading, neighbors} = useNeighbors();

  console.log(isLoading, neighbors);

  return null;
}

export default Neighbors;


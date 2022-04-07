import React, { useContext, useEffect, useMemo, useState } from "react";
import { Context } from "shared/Context";
import api from "shared/api";
import styled from "styled-components";
import Selector from "components/Selector";

import Loading from "components/Loading";

import _ from "lodash";
import DeploymentCard from "./DeploymentCard";
import { Environment, PRDeployment } from "../types";
import { useRouting } from "shared/routing";
import { useHistory, useLocation } from "react-router";

const AvailableStatusFilters = [
  "all",
  "creating",
  "failed",
  "active",
  "inactive",
];

type AvailableStatusFiltersType = typeof AvailableStatusFilters[number];

const DeploymentList = ({ environments }: { environments: Environment[] }) => {
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);
  const [deploymentList, setDeploymentList] = useState<PRDeployment[]>([]);
  const [
    statusSelectorVal,
    setStatusSelectorVal,
  ] = useState<AvailableStatusFiltersType>("all");
  const [selectedRepo, setSelectedRepo] = useState("all");

  const { currentProject, currentCluster } = useContext(Context);
  const { getQueryParam, pushQueryParams } = useRouting();
  const location = useLocation();
  const history = useHistory();

  const getPRDeploymentList = () => {
    return api.getPRDeploymentList(
      "<token>",
      {},
      {
        project_id: currentProject.id,
        cluster_id: currentCluster.id,
      }
    );
  };

  useEffect(() => {
    const selected_repo = getQueryParam("repository");

    const repo = environments.find(
      (env) => `${env.git_repo_owner}/${env.git_repo_name}` === selected_repo
    );

    if (!repo) {
      pushQueryParams({}, ["repository"]);
      return;
    }

    if (selected_repo !== selectedRepo) {
      setSelectedRepo(`${repo.git_repo_owner}/${repo.git_repo_name}`);
    }
  }, [location.search, history]);

  useEffect(() => {
    const status_filter = getQueryParam("status_filter");

    if (!AvailableStatusFilters.includes(status_filter)) {
      pushQueryParams({}, ["status_filter"]);
      return;
    }

    if (status_filter !== statusSelectorVal) {
      setStatusSelectorVal(status_filter);
    }
  }, [location.search, history]);

  useEffect(() => {
    pushQueryParams({}, ["status_filter", "repository"]);
  }, []);

  useEffect(() => {
    let isSubscribed = true;
    getPRDeploymentList()
      .then(({ data }) => {
        if (!isSubscribed) {
          return;
        }

        setDeploymentList(data.deployments || []);
        setIsLoading(false);
      })
      .catch((err) => {
        console.error(err);
        if (isSubscribed) {
          setHasError(true);
        }
      });

    return () => {
      isSubscribed = false;
    };
  }, [currentCluster, currentProject, statusSelectorVal]);

  const handleRefresh = () => {
    setIsLoading(true);
    getPRDeploymentList()
      .then(({ data }) => {
        setDeploymentList(data.deployments || []);
      })
      .catch((err) => {
        setHasError(true);
        console.error(err);
      })
      .finally(() => setIsLoading(false));
  };

  if (hasError) {
    return <Placeholder>Error</Placeholder>;
  }

  const filteredDeployments = useMemo(() => {
    return deploymentList.filter((d) => {
      return d.status === statusSelectorVal;
    });
  }, [statusSelectorVal]);

  const renderDeploymentList = () => {
    if (isLoading) {
      return (
        <Placeholder>
          <Loading />
        </Placeholder>
      );
    }

    if (!deploymentList.length) {
      return (
        <Placeholder>
          No preview apps have been found. Open a PR to create a new preview
          app.
        </Placeholder>
      );
    }

    if (!filteredDeployments.length) {
      return (
        <Placeholder>
          No preview apps have been found with the given filter.
        </Placeholder>
      );
    }

    return filteredDeployments.map((d) => {
      return <DeploymentCard deployment={d} onDelete={handleRefresh} />;
    });
  };

  const repoOptions = environments
    .map((env) => ({
      label: `${env.git_repo_owner}/${env.git_repo_name}`,
      value: `${env.git_repo_owner}/${env.git_repo_name}`,
    }))
    .concat({
      label: "All",
      value: "all",
    });

  const handleStatusFilterChange = (value: string) => {
    pushQueryParams({ status_filter: value });
    setStatusSelectorVal(value);
  };

  const handleRepoFilterChange = (value: string) => {
    pushQueryParams({ repository: value });
    setSelectedRepo(value);
  };

  return (
    <Container>
      <ControlRow>
        <ActionsWrapper>
          <StyledStatusSelector>
            <Label>
              <i className="material-icons">filter_alt</i>
              Status
            </Label>
            <Selector
              activeValue={statusSelectorVal}
              setActiveValue={handleStatusFilterChange}
              options={[
                {
                  value: "all",
                  label: "All",
                },
                {
                  value: "creating",
                  label: "Creating",
                },
                {
                  value: "failed",
                  label: "Failed",
                },
                {
                  value: "active",
                  label: "Active",
                },
                {
                  value: "inactive",
                  label: "Inactive",
                },
              ]}
              dropdownLabel="Status"
              width="150px"
              dropdownWidth="230px"
              closeOverlay={true}
            />
          </StyledStatusSelector>
          <StyledStatusSelector>
            <Label>
              <i className="material-icons">filter_alt</i>
              Repository
            </Label>
            <Selector
              activeValue={selectedRepo}
              setActiveValue={handleRepoFilterChange}
              options={repoOptions}
              dropdownLabel="Repository"
              width="200px"
              dropdownWidth="300px"
              closeOverlay
            />
          </StyledStatusSelector>

          <RefreshButton color={"#7d7d81"} onClick={handleRefresh}>
            <i className="material-icons">refresh</i>
          </RefreshButton>
        </ActionsWrapper>
      </ControlRow>
      <EventsGrid>{renderDeploymentList()}</EventsGrid>
    </Container>
  );
};

export default DeploymentList;

const ActionsWrapper = styled.div`
  display: flex;
`;

const RefreshButton = styled.button`
  display: flex;
  align-items: center;
  justify-content: center;
  color: ${(props: { color: string }) => props.color};
  cursor: pointer;
  border: none;
  background: none;
  border-radius: 50%;
  margin-left: 10px;
  > i {
    font-size: 20px;
  }
  :hover {
    background-color: rgb(97 98 102 / 44%);
    color: white;
  }
`;

const Placeholder = styled.div`
  padding: 30px;
  margin-top: 35px;
  padding-bottom: 40px;
  font-size: 13px;
  color: #ffffff44;
  min-height: 400px;
  height: 50vh;
  background: #ffffff11;
  border-radius: 8px;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;

  > i {
    font-size: 18px;
    margin-right: 8px;
  }
`;

const Container = styled.div`
  margin-top: 33px;
  padding-bottom: 120px;
`;

const ControlRow = styled.div`
  display: flex;
  margin-left: auto;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 35px;
  padding-left: 0px;
`;

const EventsGrid = styled.div`
  display: grid;
  grid-row-gap: 20px;
  grid-template-columns: 1;
`;

const StyledStatusSelector = styled.div`
  display: flex;
  align-items: center;
  font-size: 13px;
  :not(:first-child) {
    margin-left: 15px;
  }
`;

const Header = styled.div`
  font-weight: 500;
  color: #aaaabb;
  font-size: 16px;
  margin-bottom: 15px;
  width: 50%;
`;

const Subheader = styled.div`
  width: 50%;
`;

const Label = styled.div`
  display: flex;
  align-items: center;
  margin-right: 12px;

  > i {
    margin-right: 8px;
    font-size: 18px;
  }
`;

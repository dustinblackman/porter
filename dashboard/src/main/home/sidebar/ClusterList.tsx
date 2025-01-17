import React, { useState, useEffect, useRef, useContext } from "react";
import styled from "styled-components";
import gradient from "assets/gradient.png";
import api from "shared/api";
import infra from "assets/infra.png";

import { Context } from "shared/Context";
import { ClusterType } from "shared/types";
import { RouteComponentProps, withRouter } from "react-router";
import Icon from "components/porter/Icon";
import Spacer from "components/porter/Spacer";
import { pushFiltered } from "shared/routing";
import SidebarLink from "./SidebarLink";

const ClusterList: React.FC<PropsType> = (props) => {
  const { setCurrentCluster, user, currentCluster, currentProject } = useContext(Context);
  const [expanded, setExpanded] = useState<boolean>(false);
  const wrapperRef = useRef<HTMLDivElement>(null);
  const [clusters, setClusters] = useState<ClusterType[]>([]);
  const [options, setOptions] = useState<any[]>([]);

  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (
        wrapperRef.current &&
        !wrapperRef.current.contains(e.target as Node)
      ) {
        setExpanded(false);
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);

  useEffect(() => {
    if (currentProject) {
      api
        .getClusters("<token>", {}, { id: currentProject?.id })
        .then((res) => {
          if (res.data) {
            let clusters = res.data;
            clusters.sort((a: any, b: any) => a.id - b.id);
            if (clusters.length > 0) {
              let options = clusters.map((item: { name: any; vanity_name: string; }) => ({
                label: (item.vanity_name ? item.vanity_name : item.name),
                value: item.name
              }));
              setClusters(clusters);
              setOptions(options);
            }
          }
        });
    }
  }, [currentProject]);
  const truncate = (input: string) => input.length > 21 ? `${input.substring(0, 21)}...` : input;

  const renderOptionList = () =>
    options.map((option, i: number) => (
      <Option
        key={i}
        selected={option.value === currentCluster?.name}
        onClick={() => {
          setExpanded(false);
          const cluster = clusters.find(c => c.name === option.value);
          setCurrentCluster(cluster);
          pushFiltered(props, "/apps", ["project_id"], {});
        }}
      >

        <Icon src={infra} height={"14px"} />
        <ClusterLabel>{option.label}</ClusterLabel>
      </Option>
    ));

  const renderDropdown = () =>
    expanded && (
      <Dropdown>
        {renderOptionList()}
      </Dropdown>
    );

  if (currentCluster) {
    return (
      <StyledClusterSection ref={wrapperRef}>
        <MainSelector
          onClick={() => setExpanded(!expanded)}
          expanded={expanded}
        >

          <NavButton>
            <Img src={infra} />
            <ClusterName>{truncate(currentCluster.vanity_name ? currentCluster.vanity_name : currentCluster?.name)}</ClusterName>

            {clusters.length > 1 && <i className="material-icons">arrow_drop_down</i>}
          </NavButton>
        </MainSelector>
        {clusters.length > 1 && renderDropdown()}
      </StyledClusterSection >
    );
  }

  return (
    <InitializeButton
      onClick={() =>
        pushFiltered(props, "/new-cluster", ["cluster_id"], {
          new_cluster: true,
        })
      }
    >
      <Plus>+</Plus> Create a cluster
    </InitializeButton>
  );
};

export default withRouter(ClusterList);

const ClusterLabel = styled.div`
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  flex-grow: 1;
  margin-left: 5px
`;

const Plus = styled.div`
  margin-right: 10px;
  font-size: 15px;
`;

const InitializeButton = styled.div`
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: calc(100% - 30px);
  height: 38px;
  margin: 8px 15px;
  font-size: 13px;
  font-weight: 500;
  border-radius: 3px;
  color: ${props => props.theme.text.primary};
  padding-bottom: 1px;
  cursor: pointer;
  background: #ffffff11;

  :hover {
    background: #ffffff22;
  }
`;

const Option = styled.div`
  width: 100%;
  border-top: 1px solid #00000000;
  border-bottom: 1px solid
    ${(props: { selected: boolean; lastItem?: boolean }) =>
    props.lastItem ? "#ffffff00" : "#ffffff15"};
  height: 45px;
  display: flex;
  align-items: center;
  font-size: 13px;
  align-items: center;
  padding-left: 10px;
  cursor: pointer;
  padding-right: 10px;
  background: ${(props: { selected: boolean; lastItem?: boolean }) =>
    props.selected ? "#ffffff11" : ""};
  :hover {
    background: ${(props: { selected: boolean; lastItem?: boolean }) =>
    props.selected ? "" : "#ffffff22"};
  }

  > i {
    font-size: 18px;
    margin-right: 12px;
    margin-left: 5px;
    color: #ffffff44;
  }
`;

const Dropdown = styled.div`
  position: absolute;
  right: 13px;
  top: calc(100% + 5px);
  background: #171b20;
  width: 210px;
  max-height: 500px;
  border-radius: 3px;
  z-index: 999;
  overflow-y: auto;
  margin-bottom: 20px;
  box-shadow: 0 5px 15px 5px #00000077;
`;

const ClusterName = styled.div`
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  display: flex;
  align-items: center;
  max-width: 180px; // You can adjust this value according to your needs
`;

const MainSelector = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;

  font-size: 14px;
  cursor: pointer;
  padding: 10px 0;
  
  :hover {
    > i {
      background: #ffffff22;
    }
  }

  > i {
    margin-left: 0px;
    margin-right: 0px;
    font-size: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 20px;
    background: ${(props: { expanded: boolean }) =>
    props.expanded ? "#ffffff22" : ""};
  }
`;

const StyledClusterSection = styled.div`
  position: relative;
  margin-left: 3px;
  background: #545ec7;
  border-right: 1px solid #2c2e31;
`;
const NavButton = styled(SidebarLink)`
  display: flex;
  align-items: center;
  position: relative;
  text-decoration: none;
  border-radius: 5px;
  margin-left: 16px;
  font-size: 13px;
  color: ${props => props.theme.text.primary};
  cursor: ${(props: { disabled?: boolean }) =>
    props.disabled ? "not-allowed" : "pointer"};

  background: ${(props: any) => (props.active ? "#ffffff11" : "")};

  :hover {
    background: ${(props: any) => (props.active ? "#ffffff11" : "#ffffff08")};
  }

  &.active {
    background: #ffffff11;

    :hover {
      background: #ffffff11;
    }
  }

  :hover {
    background: #ffffff08;
  }

  > i {
    font-size: 18px;
    border-radius: 0px;
    margin-left: 2px;
    margin-right: 0px;
  }
`;


const Img = styled.img<{ enlarge?: boolean }>`
  padding: ${(props) => (props.enlarge ? "0 0 0 1px" : "4px")};
  height: 22px;
  padding-top: 4px;
  border-radius: 3px;
  margin-right: 8px;
  margin-left: 2px;
  opacity: 0.8;
`;

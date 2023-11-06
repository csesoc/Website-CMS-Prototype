# Backend

The backend folder contains all our backend code 🤯. Theres a few important folders:
 - ### `database/`
   - Package contains all database specific code (basically just contexts + repositories)
   - Repositories are repositories, just provide methods for interacting with the database
   - Contexts are actual database connections, theres a `testing_context` and a normal `live_context`, the `testing_context` wraps all SQL queries in a transaction and rolls them back once testing has finished and `live_context` does nothing special 🙁
 - ### `endpoints/`
   - Contains all our HTTP handlers + methods for decorating those handlers, additionally provides methods for attaching handlers to a `http.ServeMux`
 - ### `editor/`
   - There are currently 3 different editor backends, once the OT backend is fully complete this will collapse down to just OT
      - The OT folder contains our implementation of the operational transform algorithm, specifically Google WAVE OT
      - The pessimistic editor is a lock based editor, instead of implementing an optimistic concurrency protocol it simply locks the editor while a conflicting client is using it
      - The diffsync folder contains all the backend code for synchronisation (bit of a mess right now 😭). It is essentially just a big implementation for the differential sync algorithm and supports extension via the construction of `extensions`
   - The package allows multiple clients to be connected to the same document and ensures that each client sees the exact same document state (via a best effort diff + patch)
   - Currently only works on strings with no inheret structure but in the future will work on JSON documents too 😀 
 - ### `environment/` & `internal/`
   - Methods/packages that can be used to retrieve information about the current environment and other internal utilities 
 - ### `client`
   - WIP TypeScript implementation of client server for operational transform


## Architectural Overview
At a very high level our CMS looks as follows
![CMS arch](./docs/CMS%20high%20level%20architecture.png)
There are a few core parts of the CMS
 - Handlers (within the `endpoints` directory)
 - Repository Layer (within the `database` directory)
 - Concurrent Editor (within the `editor` directory)

Each of these sections has their own little bit of documentation in their respective directory so this bit will focus more on the high level aspects of the CMS.

At a very high level the CMS allows users to easily manage, collaborate on and create static web content, the whole goal of this is to allow for the easy extension of maintenance of the CSESoc Website as well as any other websites CSESoc produces. The CMS identifies different websites as "frontends", each frontend gets its own unique file tree that members of the frontend can manage, this file tree is stored within the `Postgres DB`. 

Users of the CMS work on the idea of "documents", documents are unique pages that they can create and edit, there are two different types of documents: `published` and `unpublished`, published documents are visible to any unauthenticated user, regardless of if they are a member of the frontend group or not while `unpublished` documents are only visible to editors. Each of these document types are stored in their respective content volumes, we duplicate these documents as we may have a version A of a document that we wish to be public but have a version A' that we're still editing, that should be private. We also store other types of content such as images and videos but such content does not receive a published/unpublished distinction.

## Papers worth Reading
Most of the complexity of the CMS backend is within the editor, to aid with your tickets we have accumilated a few great resources that are worth a read.
  - [A survey of OT algorithms](https://www.researchgate.net/profile/Ajay-Khunteta-2/publication/45183356_A_Survey_on_Operational_Transformation_Algorithms_Challenges_Issues_and_Achievements/links/5b9b27dca6fdccd3cb533171/A-Survey-on-Operational-Transformation-Algorithms-Challenges-Issues-and-Achievements.pdf?origin=publication_detail)
  - [The Jupiter Operational Transform Algorithm](https://lively-kernel.org/repository/webwerkstatt/projects/Collaboration/paper/Jupiter.pdf)
  - [Google Wave OT (Multi-client single server OT)](https://svn.apache.org/repos/asf/incubator/wave/whitepapers/operational-transform/operational-transform.html)
  - [Transformation algorithms for ordered n-ary trees](https://arxiv.org/pdf/1512.05949.pdf)
  - [Differental Synchronisation](https://neil.fraser.name/writing/sync/eng047-fraser.pdf)

## Language Documentation & Resources
If this is your first time using Go then the following resources might be of some use:
 - [2022 CSESoc Dev Go Workshop](https://drive.google.com/file/d/1zLJHkcktLFXKXs6MFNVmWyfVtTHq8ng8/view)
 - [Unit testing in Go](https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package)
    - [Interface mocking & gomock](https://itnext.io/how-to-write-better-unit-tests-in-go-using-mocks-4dd05e867b17)
 - [Generic Programming in Go](https://go.dev/doc/tutorial/generics)
 - [Concurrent Programming in Go](https://golangdocs.com/concurrency-in-golang)
      - [Communicating Sequential Proceseses (Go's model of concurrency)](https://www.cs.cmu.edu/~crary/819-f09/Hoare78.pdf)
          - It's not required but its a nice read
 - [Godoc](https://go.dev/blog/godoc)
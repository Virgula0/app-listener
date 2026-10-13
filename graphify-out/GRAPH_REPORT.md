# Graph Report - app-listener  (2026-10-13)

## Corpus Check
- 365 files · ~1,561,915 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4689 nodes · 19066 edges · 208 communities (179 shown, 29 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1365 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1107d9c9`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- orphans_test.go
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- go_pkg_path_filepath
- planUpdaters
- TestMonitorUseCaseLifecycle
- Candidate
- install.go
- wipe
- .startGuardStd
- guard_trust.bpf.c
- SetSectionWhitelist
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- Ledger
- hashFileShared
- Monitor
- serve.go
- BinaryEntry
- IntegrationSuite
- fileedit_model_test.go
- IntegrationSuite
- headerWidget
- backups.go
- controlServer
- bpfstats/main.go
- Finder
- New
- IntegrationSuite
- liveSession
- IntegrationSuite
- update.go
- highlight_test.go
- .syncSetsLocked
- TrustGuard
- tempRow
- inode_dev
- main_test.go
- daemonconfig.go
- github.com/charmbracelet/bubbletea.Cmd
- procstats.go
- codeedit.go
- sync.Mutex
- engine
- GuardEvent
- audit.go
- .unlockFileInPlace
- exe_supersede.h
- fakeGuardRepo
- errno
- pinstate.go
- Vault
- buildTrustedSet
- testcontainers.Container
- fcntl
- copyRegularAt
- runNetworkGuard
- watchSet
- btrfsLayoutFrom
- GenerateConf
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- process_vm_readv.c
- OpenSystemPlaced
- downloadAndVerify
- model
- newTestVetter
- GuardInodeKey
- ResolvePinBase
- github.com/charmbracelet/bubbles/textarea.Model
- update_test.go
- sealFileVaultWithMasterKey
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- safeio.go
- .judgeProcess
- IntegrationSuite
- Well-known bypass classes
- runNetworkMonitor
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- CandidateDir
- Build & sign release assets (reusable workflow)
- runMonitor
- fileedit_symlink_test.go
- net_tester/main.go
- NewEventFanout
- app-listener Terminal Demo Recording
- IntegrationSuite
- time.Duration
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- catalogRefresher
- NewDumpHook
- WebSocket /ws client
- netModel
- string
- install.sh
- mc_tag.h
- newChangelogModel
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- Run
- trace-app-libs.sh
- IntegrationSuite
- serveWebSocket
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- DiscoverForUsers
- Config
- preload_wait.c
- Guard
- go_pkg_fmt
- runUpdate
- Vault
- ListUsers
- StripPinnedTempAllows
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- shQuote
- inodeChain
- Highlighter
- KernelDev
- Installation
- os.File
- graphify-refresh.sh
- New
- bufio.Reader
- github.com/cilium/ebpf/link.Link
- BinaryStat
- IntegrationSuite
- is_event_type_allowed
- bottomBar
- lockAndDeprovision
- parseGuardEvent
- RefreshCatalog
- Limitations
- Troubleshooting
- time.Time
- openBunTmp
- Daemon configuration
- edit-protected
- absPath
- runDaemon
- Compatibility
- How it works
- ResolveTempBinary
- loadDaemonConfig
- bpfLsmListed
- launch_rule_at
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- watchPattern
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- integrationtests/guard_test.go
- SanitizeText
- findTopLevelWindow
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- noticeModel
- hashCache
- assetsFor
- io.Reader
- .TestGuard_PathRmdirEvictsInodeImmediately
- TestBatchInodeRows
- newTestSession

## God Nodes (most connected - your core abstractions)
1. `Config` - 95 edges
2. `Guard` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
9. `fileEditModel` - 60 edges
10. `User` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (208 total, 29 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.16
Nodes (24): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, multicallCacheKey, launchModel() (+16 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (132): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+124 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.12
Nodes (6): IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), daemonEvent

### Community 3 - "orphans_test.go"
Cohesion: 0.21
Nodes (17): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors() (+9 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (77): TestSetSectionWhitelistPreservesGroupStructure(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+69 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (38): admit_refusal(), check_watched(), code_vouched(), emit_gate(), env_loader(), env_step(), exe_inode_of(), fill_current_image() (+30 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.09
Nodes (58): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), catalogPatterns(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor() (+50 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (15): addErrKind, guardUnitTest, classifyAddErr(), eventMask(), guardModeKey(), addRecorder(), copySelf(), runTool() (+7 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 11 - "planUpdaters"
Cohesion: 0.16
Nodes (16): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+8 more)

### Community 12 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.24
Nodes (11): NewGuardUseCase(), NewNetworkMonitorUseCase(), newFakeMonitorRepo(), newFakeNetworkGuardRepo(), newFakeNetworkMonitorRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle() (+3 more)

### Community 13 - "Candidate"
Cohesion: 0.21
Nodes (16): appendSectionsAndEdit(), collectDiffAdditions(), configPaths(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig() (+8 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (56): reloadDaemonForPasswordChange(), applyDiffAdditions(), restoreDaemonAfterDiffAbort(), runDiffCatalog(), writeEditPasswordHash(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled() (+48 more)

### Community 15 - "wipe"
Cohesion: 0.20
Nodes (10): isFileVaultCiphertext(), isRegularFileTarget(), hasEncryptionPolicy(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall(), TestNewBoundedKeyFnRetryAborts() (+2 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.11
Nodes (51): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+43 more)

### Community 18 - "SetSectionWhitelist"
Cohesion: 0.07
Nodes (37): logRefreshChange(), TestLogRefreshChange(), LibraryBlock, RefreshOptions, SectionScan, ConfSafePath(), findSectionEnd(), findSectionStart() (+29 more)

### Community 19 - "fileEditModel"
Cohesion: 0.12
Nodes (9): clipLabel(), fileModeToUnixMode(), humanSize(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.31
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.10
Nodes (56): Vault, NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged() (+48 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (11): ownerPaths(), fakeDaemon, Resource, GuardRepository, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, uniqueEncryptionRoots() (+3 more)

### Community 23 - "auth.go"
Cohesion: 0.15
Nodes (27): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+19 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.11
Nodes (24): buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), selfGuards, selfProtectSpec, AdmissionCheck (+16 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.12
Nodes (51): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), count_degrade(), discover_guarded_parent(), emit_process_denial(), evict_inode_from_guard(), get_dentry_from_path() (+43 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.11
Nodes (10): NetBpfEvent, NetEvent, FormatAddr(), Ntohs(), NewNetworkMonitor(), statInodeKey(), TestNetworkMonitor_NewFailure(), TestStatInodeKey() (+2 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "Ledger"
Cohesion: 0.12
Nodes (10): candidates(), confirmable(), confirmAll(), describe(), hashNow(), run(), Ledger, Pending (+2 more)

### Community 29 - "hashFileShared"
Cohesion: 0.29
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 30 - "Monitor"
Cohesion: 0.06
Nodes (8): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), dirTarget(), Monitor, monitorUnitTest, pathCache, pathCacheEntry

### Community 31 - "serve.go"
Cohesion: 0.11
Nodes (22): basicAuth(), decodeEndsAtEOF(), isInteractiveTerminal(), newServeListener(), newServeServer(), parseResize(), readResizeLoop(), relayServeSignals() (+14 more)

### Community 32 - "BinaryEntry"
Cohesion: 0.09
Nodes (18): deferredBinary, BinaryEntry, BinariesSummary(), canonicalBinaryPath(), canonicalPaths(), confinedEntry(), patternMatches(), readPinnedResID() (+10 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.17
Nodes (24): isOctalMode(), newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow() (+16 more)

### Community 35 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 36 - "headerWidget"
Cohesion: 0.16
Nodes (3): dataRowRenderer, headerRenderer, headerWidget

### Community 37 - "backups.go"
Cohesion: 0.13
Nodes (23): deletePostBackups(), restoreBackups(), pickUsers(), offerBackupCleanup(), runUninstall(), Delete(), Find(), find() (+15 more)

### Community 38 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 39 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 40 - "Finder"
Cohesion: 0.17
Nodes (9): findButton, span, column(), findAll(), Finder, lowerRunes(), overlay(), runesEqual() (+1 more)

### Community 41 - "New"
Cohesion: 0.22
Nodes (15): harnessSuite, looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+7 more)

### Community 44 - "IntegrationSuite"
Cohesion: 0.18
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 45 - "update.go"
Cohesion: 0.16
Nodes (10): compareStableVersions(), parseChecksum(), parsePublicKey(), parseStableVersion(), parseVerPart(), TestParseChecksum(), TestParseStableVersion(), verifyRelease() (+2 more)

### Community 46 - "highlight_test.go"
Cohesion: 0.20
Nodes (14): NewHighlighter(), editSession(), sampleGo(), TestDaemonConfTokens(), TestGutterKeepsTextColumnAligned(), TestHighlighterColorsTokens(), TestHighlighterFollowsCursorAlone(), TestHighlighterOnlyChangesColors() (+6 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "TrustGuard"
Cohesion: 0.07
Nodes (15): trustHook, trustRows, deleteResKeys(), setSupersedeTrust(), syncU32Map(), cStr(), syncMap(), TrustGuard (+7 more)

### Community 50 - "tempRow"
Cohesion: 0.18
Nodes (6): exeRow, tempGrant, tempMaskOp, tempRow, Guard, liveTrust()

### Community 51 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), code_suspect(), current_is_inspector() (+15 more)

### Community 52 - "main_test.go"
Cohesion: 0.05
Nodes (20): IntegrationSuite, decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest (+12 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.16
Nodes (33): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+25 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.17
Nodes (4): newFindInput(), clamp(), serveTestModel, sessionModel

### Community 55 - "procstats.go"
Cohesion: 0.18
Nodes (9): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+1 more)

### Community 56 - "codeedit.go"
Cohesion: 0.25
Nodes (16): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), Navigate(), onlySpaces() (+8 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 58 - "engine"
Cohesion: 0.07
Nodes (12): TestNearestRootClaims(), TestTieRank(), engine, nearestRootClaims(), pathWithin(), TestPathWithinMatchesPrefixRule(), tieRank(), deleteInoKeys() (+4 more)

### Community 59 - "GuardEvent"
Cohesion: 0.07
Nodes (27): displayEntries(), holdForward(), holdHeadless(), pingInterval(), chooseResource(), runLiveEdit(), HeadlessLine(), resolveGuardConfig() (+19 more)

### Community 60 - "audit.go"
Cohesion: 0.23
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 61 - ".unlockFileInPlace"
Cohesion: 0.25
Nodes (15): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, openFileVaultWithMasterKey(), recoverFileInPlace(), rejectSidecarSymlink() (+7 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.17
Nodes (18): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+10 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "pinstate.go"
Cohesion: 0.13
Nodes (28): configureDaemonLogging(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive() (+20 more)

### Community 66 - "Vault"
Cohesion: 0.23
Nodes (7): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), Vault, readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 68 - "buildTrustedSet"
Cohesion: 0.21
Nodes (13): inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree(), keys(), TestWarnUntrustedLibs_OnlyWhatTheKernelRefuses() (+5 more)

### Community 69 - "testcontainers.Container"
Cohesion: 0.19
Nodes (3): configBinaryLines(), rawExec(), watchPathsInConfig()

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "runNetworkGuard"
Cohesion: 0.11
Nodes (17): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, NetGuardEvent (+9 more)

### Community 73 - "watchSet"
Cohesion: 0.33
Nodes (4): isSymlink(), watchDir, watchPos, watchSet

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 75 - "GenerateConf"
Cohesion: 0.11
Nodes (28): TestDiffMergePreservesExistingAndParses(), sectionsFromCandidates(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestGeneratedLibraryBlockRoundTrips() (+20 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (19): AddServeFlags(), credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+11 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "OpenSystemPlaced"
Cohesion: 0.09
Nodes (21): setConfinedHomes(), placedWalk, TrustGuard, launchKey(), TestLaunchRulesFitTheKey(), confineBelow(), inUserHome(), SetUserHomes() (+13 more)

### Community 81 - "downloadAndVerify"
Cohesion: 0.20
Nodes (10): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestSanityCheckBinaryRefusesOtherMachine(), updateHTTPClient() (+2 more)

### Community 82 - "model"
Cohesion: 0.22
Nodes (7): formatType(), listenForEvents(), syncViewport(), tickStats(), eventLine, eventMsg, model

### Community 83 - "newTestVetter"
Cohesion: 0.18
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 84 - "GuardInodeKey"
Cohesion: 0.05
Nodes (31): refusedReplacements, checkBtrfsKeys(), ensureBtrfsLayout(), appendInodeRow(), batchInodeRows(), inodeRowsOf(), engine, Inspector (+23 more)

### Community 85 - "ResolvePinBase"
Cohesion: 0.50
Nodes (5): dirIsRootOwnedSafe(), isBpffs(), mountBpffs(), ResolvePinBase(), TestDirIsRootOwnedSafe()

### Community 86 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.16
Nodes (14): lineStyles, viewportState, fg(), baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle() (+6 more)

### Community 87 - "update_test.go"
Cohesion: 0.12
Nodes (12): fetchReleases(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel(), testKeyPair(), TestNewerThanInstalled() (+4 more)

### Community 88 - "sealFileVaultWithMasterKey"
Cohesion: 0.26
Nodes (11): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape() (+3 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.07
Nodes (30): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), ldConfInclude(), ldSoConfDirs() (+22 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.09
Nodes (36): UIDResolver, NewUIDResolver(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless(), runServedTUI(), runTUI() (+28 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.12
Nodes (10): NetEventType, FirstOnBtrfs(), NetEventTypes(), ParseNetEventType(), eventsetSummary(), eventTypeKey(), NetGuard, mcNameKey() (+2 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "safeio.go"
Cohesion: 0.07
Nodes (41): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+33 more)

### Community 100 - ".judgeProcess"
Cohesion: 0.14
Nodes (16): reasonText(), codeVouched(), NetGuard, isTraced(), judgePasses(), loaderEnv(), parseProcs(), roFor() (+8 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "runNetworkMonitor"
Cohesion: 0.12
Nodes (9): ParseNetEventsFlag(), runHeadless(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString(), NetworkMonitorRepository, NewNetModel() (+1 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "CandidateDir"
Cohesion: 0.06
Nodes (33): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+25 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.07
Nodes (27): BuildEBPFTargets(), CheckEBPF(), RawTarget, isSubDir(), MakeDisplayPaths(), ParseEventsFlag(), pathWithinDir(), ResolveTargets() (+19 more)

### Community 110 - "fileedit_symlink_test.go"
Cohesion: 0.23
Nodes (11): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+3 more)

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "NewEventFanout"
Cohesion: 0.24
Nodes (6): NewEventFanout(), TestEventFanoutDropOldestUnderPressure(), TestEventFanoutDuplicatesInOrderAndCloses(), TestEventFanoutStopClosesOutputs(), EventFanout, EventFanout[T]

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "time.Duration"
Cohesion: 0.15
Nodes (8): debouncer, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "catalogRefresher"
Cohesion: 0.14
Nodes (10): drain(), readEvents(), sameConfig(), swapFile(), grantNewBinaries(), catalogRefresher, Lines(), TestLinesKeyOnLinkAndIncludePending() (+2 more)

### Community 121 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "netModel"
Cohesion: 0.12
Nodes (13): sanitizeTerminalText(), sanitizeTerminalTexts(), clipToWidth(), formatResourceBar(), formatAddr(), formatNetProto(), formatNetType(), listenForNetEvents() (+5 more)

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "newChangelogModel"
Cohesion: 0.27
Nodes (6): newChangelogModel(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestNewChangelogModelViewContainsTitleAndNotes(), changelogModel

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.10
Nodes (12): BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock(), Run() (+4 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "Run"
Cohesion: 0.19
Nodes (7): columnState, dataRow, guiModel, newColumnState(), newDataRow(), newHeaderWidget(), Run()

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "serveWebSocket"
Cohesion: 0.16
Nodes (8): configureWebSocket(), sameOrigin(), serveWebSocket(), TestSameOrigin(), writeSnapshot(), writeSnapshotLoop(), servedModel, snapshotHub

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.13
Nodes (16): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant() (+8 more)

### Community 143 - "DiscoverForUsers"
Cohesion: 0.33
Nodes (7): Discover(), DiscoverForUsers(), DiscoverSystem(), TestDiscoverForUsersDeduplicates(), TestDiscoverOnlyExisting(), TestDiscoverPerUserHomes(), TestDiscoverSystemOnlyAbsolute()

### Community 144 - "Config"
Cohesion: 0.22
Nodes (12): newCatalogRefresher(), reloadSlotsNeeded(), resolveInspectors(), TestResolveInspectorsKeepsOnlyAdmitted(), trustManager, startTrustGuard(), trustStartupError(), inspectorAdmitter (+4 more)

### Community 146 - "Guard"
Cohesion: 0.08
Nodes (10): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), binaryVerifyState, rootHandle, openRoot(), SharedPinDegraded() (+2 more)

### Community 147 - "go_pkg_fmt"
Cohesion: 0.06
Nodes (36): TestPickEncryptedDirSingleSkipsPicker(), pickEncryptedDir(), runOfflineEdit(), atomicWriteAt(), descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent() (+28 more)

### Community 148 - "runUpdate"
Cohesion: 0.21
Nodes (15): changelogText(), showChangelog(), TestChangelogText(), confirmUpdate(), filterChannel(), isTerminal(), latestTime(), newerThanInstalled() (+7 more)

### Community 149 - "Vault"
Cohesion: 0.18
Nodes (14): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+6 more)

### Community 150 - "ListUsers"
Cohesion: 0.10
Nodes (22): pathCovered(), TestPathCovered(), cleanOrphanedFscrypt(), isInsidePath(), askSSHAgentUsers(), askDecrypt(), cleanOrphanedMetadata(), decryptStep() (+14 more)

### Community 151 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.27
Nodes (20): ConfinedEntry(), HashStats(), newVerifyState(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash() (+12 more)

### Community 154 - "shQuote"
Cohesion: 0.20
Nodes (3): IntegrationSuite, shQuote(), IntegrationSuite

### Community 155 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 156 - "Highlighter"
Cohesion: 0.17
Nodes (6): editorModel, Highlighter, lexerFor(), shebangLexer(), tokenClassOf(), newEditorModel()

### Community 157 - "KernelDev"
Cohesion: 0.16
Nodes (16): fdStat, physicalParent(), statFD(), isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed() (+8 more)

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.05
Nodes (55): ReplacementCheck, VettedInode, Multicall, MulticallKind, TrustGuard, ExeKey(), Guard, IsMulticallRefusal() (+47 more)

### Community 161 - "New"
Cohesion: 0.26
Nodes (11): diffModel, moveTo(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+3 more)

### Community 162 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 163 - "github.com/cilium/ebpf/link.Link"
Cohesion: 0.33
Nodes (3): foreignPinnedLink(), NetGuard, hook

### Community 164 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 166 - "is_event_type_allowed"
Cohesion: 0.44
Nodes (10): emit_event(), guard_net_socket_bind(), guard_net_socket_connect(), guard_net_socket_listen(), guard_net_socket_recvmsg(), guard_net_socket_sendmsg(), is_event_type_allowed(), is_socket_guarded() (+2 more)

### Community 167 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 168 - "lockAndDeprovision"
Cohesion: 0.19
Nodes (10): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision(), CopyTreeWithProgress() (+2 more)

### Community 169 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 170 - "RefreshCatalog"
Cohesion: 0.40
Nodes (5): RefreshCatalog(), TestRefreshCatalogAdmitVetsOnlyNewLines(), TestRefreshCatalogKeepsSymlinkLineNotTarget(), TestRefreshCatalogReportsVersionBump(), TestRefreshCatalogUserSectionsOnly()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 174 - "time.Time"
Cohesion: 0.31
Nodes (5): newWatchSet(), fakeInotify, sweepID, sweepMemo, Guard

### Community 175 - "openBunTmp"
Cohesion: 0.39
Nodes (8): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), bunTmp, inodeID

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "absPath"
Cohesion: 0.13
Nodes (5): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath()

### Community 179 - "runDaemon"
Cohesion: 0.09
Nodes (31): CheckBPFLSM(), backingDeviceUnion(), buildConcurrency(), buildGuards(), catchLifecycleSignals(), confirmMasterKeyOverwrite(), makeReloadHandler(), newPinGeneration() (+23 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "ResolveTempBinary"
Cohesion: 0.22
Nodes (9): chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward(), binaryStatOf(), ResolveTempBinary() (+1 more)

### Community 183 - "loadDaemonConfig"
Cohesion: 0.15
Nodes (12): loadDaemonConfig(), markReady(), noteReloadedSettings(), notifySystemdReady(), resolveConfigPath(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), engine (+4 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 185 - "launch_rule_at"
Cohesion: 0.40
Nodes (3): launch_rule_at(), launch_step(), launch_token()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 189 - "watchPattern"
Cohesion: 0.43
Nodes (5): addSystemBinary(), catalogWatchPlan(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), watchPattern

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 194 - "integrationtests/guard_test.go"
Cohesion: 0.32
Nodes (5): attrOpCase, guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 195 - "SanitizeText"
Cohesion: 0.11
Nodes (16): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), systemPatterns(), systemRule(), binaryVetter, vetDecision, Line() (+8 more)

### Community 196 - "findTopLevelWindow"
Cohesion: 0.38
Nodes (4): findTopLevelWindow(), floatWindow(), internAtom(), setFloatHint()

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.12
Nodes (17): checkFreshDir(), openDirNoFollow(), renewDir(), TestRuntimeClassMarkers(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), applyNewFileMeta(), defaultNewFileMode() (+9 more)

### Community 200 - "noticeModel"
Cohesion: 0.38
Nodes (3): renderBullets(), TestRenderBulletsVerbatimAndIndented(), noticeModel

### Community 201 - "hashCache"
Cohesion: 0.18
Nodes (8): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache, Ledger

### Community 202 - "assetsFor"
Cohesion: 0.33
Nodes (6): assetsFor(), checkUpdatePreconditions(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestValidateChannel(), validateChannel()

### Community 203 - "io.Reader"
Cohesion: 0.33
Nodes (5): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 205 - "TestBatchInodeRows"
Cohesion: 0.50
Nodes (5): inodeBatchFn, fakeBatches(), TestBatchInodeRows(), countInodeRows(), TestCountInodeRows()

### Community 206 - "newTestSession"
Cohesion: 0.50
Nodes (3): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity()

## Knowledge Gaps
- **147 isolated node(s):** `graphify-refresh.sh script`, `trustManager`, `controlManager`, `controlServer`, `tempJournalRow` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 485 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **29 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `testcontainers.Container`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.startGuardStd`, `main_test.go`, `IntegrationSuite`, `IntegrationSuite`, `shQuote`?**
  _High betweenness centrality (0.049) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `BinaryEntry`, `SanitizeText`, `time.Time`, `runDaemon`, `time.Duration`, `GuardInodeKey`, `Resource`, `buildOneGuard`, `sync.Mutex`, `GuardEvent`, `DaemonEvent`, `os.File`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Why does `fileEditModel` connect `fileEditModel` to `fileedit_model_test.go`, `vaultFS`, `Finder`, `go_pkg_path_filepath`, `fileedit_symlink_test.go`, `highlight_test.go`, `github.com/charmbracelet/bubbletea.Cmd`, `github.com/charmbracelet/bubbles/textarea.Model`, `Highlighter`?**
  _High betweenness centrality (0.027) - this node is a cross-community bridge._
- **What connects `graphify-refresh.sh script`, `trustManager`, `controlManager` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.022112494632889654 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.11627906976744186 - nodes in this community are weakly interconnected._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07092907092907093 - nodes in this community are weakly interconnected._
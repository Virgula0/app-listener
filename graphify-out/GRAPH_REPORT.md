# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 536 files · ~1,760,663 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 6997 nodes · 31116 edges · 230 communities (195 shown, 35 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2886 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `11775de2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- Vault
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- backend.go
- New
- protectedPage
- os.File
- plan_test.go
- github.com/charmbracelet/bubbletea.Cmd
- IntegrationSuite
- guard_trust.bpf.c
- fyne.io/fyne/v2/widget.Label
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- Config
- guard.bpf.c
- FileTarget
- find_test.go
- KernelDev
- golang.org/x/sys/unix.Stat_t
- User
- TestMonitorUseCaseLifecycle
- policyEditor
- Guard
- picker
- SelectOrphans
- IntegrationSuite
- Message
- controlServer
- render.go
- pinstate_test.go
- IntegrationSuite
- absPath
- monitorUnitTest
- install_ui.go
- fileedit_model_test.go
- EventType
- stdio
- update/update.go
- fadeScroll
- vaultfs_test.go
- editView
- Client
- daemonconfig.go
- SanitizeText
- .apply
- go_pkg_strings
- mutedText
- engine
- sanitizeTerminalText
- audit.go
- Recorder
- inode_dev
- systemd.go
- installPage
- unix_connect.c
- wm_test.go
- IntegrationSuite
- linkView
- time.Duration
- CandidateDir
- RefreshSection
- Candidate
- fcntl
- guardSpec
- planUpdaters
- codecheck/main.go
- fakeGuardRepo
- Network
- process_vm_readv.c
- startCatalogRefresh
- op
- LibraryClosure
- rcUser
- .FilterExistingWhitelist
- newTestTarget
- runForward
- updatePage
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- Run
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- Pulse
- BinaryEntry
- lockAndDeprovision
- Well-known bypass classes
- recoverFileInPlace
- jit_provenance.c
- inspector_probe.c
- recordops.c
- Vault
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- openFileVaultWithMasterKey
- cursor.go
- github.com/charmbracelet/bubbles/textarea.Model
- recordFlow
- IntegrationSuite
- go_pkg_path_filepath
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- .TestRecord_GuardParity
- runLockdown
- WebSocket /ws client
- runMonitor
- install.sh
- mc_tag.h
- Monitor
- make test-integration (rootful Docker suite)
- IntegrationSuite
- dialog_test.go
- Open
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- catalogRefresher
- serve.go
- NetEventType
- secureResources
- Ledger
- runOfflineEdit
- GlobReservations
- GuardInodeKey
- bpfstats/main.go
- fakeStore
- Finder
- IntegrationSuite
- buildOneGuard
- runNetworkMonitor
- CopyTreeWithProgress
- gui/session.go
- Installation
- backend_test.go
- graphify-refresh.sh
- RawTarget
- uninstallAll
- stripTempJournal
- NetGuardEvent
- feed.go
- shell
- desktop_test.go
- watchSet
- AddServeFlags
- WithSuspectsFrom
- Limitations
- net_tester/main.go
- IntegrationSuite
- sync.Mutex
- safeio.go
- tempGrantSetup
- configedit_test.go
- loadDaemonConfig
- Compatibility
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- Binary
- fakeDaemon
- edit-protected
- NewDumpHook
- kernel-prof.sh
- Plan
- NewMonitor
- launchKey
- DaemonEvent
- Daemon configuration
- formatEntry
- ldConfParser
- Desktop GUI
- IntegrationSuite
- stripImmutableFlags
- TestCleanRel
- .eventTypeAllowed
- UI
- netinject.c
- .TestDaemon_Btrfs_CopyIoctlGate
- .awaitDaemonUpOrDead
- prctl_gate
- VouchedDevs
- TestMain
- IntegrationSuite
- GuardEvent
- supersede_probe.c
- Prereq
- mountTable
- github.com/Virgula0/app-listener/tools/codecheck
- IntegrationSuite
- Troubleshooting
- ptrace_seize.c
- .runEditorHarness
- monitor: see what touches a path

## God Nodes (most connected - your core abstractions)
1. `Config` - 127 edges
2. `Guard` - 124 edges
3. `GuardInodeKey` - 103 edges
4. `shell` - 97 edges
5. `Load()` - 93 edges
6. `absPath()` - 91 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 72 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `runConfigEditor()` --calls--> `EditText()`  [INFERRED]
  cmd/functions/install/steps.go → internal/install/editor.go

## Import Cycles
- None detected.

## Communities (230 total, 35 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.19
Nodes (22): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+14 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (145): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), TestDaemonSelfBaselineEventsNeverWrites(), TestGroupUnlockSelfEventsIsBaselinePlusMknod(), TestForwardFlagSurface(), TestSubcommandsRegistered(), TestPreVersionPattern() (+137 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (9): configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig() (+1 more)

### Community 3 - "Vault"
Cohesion: 0.15
Nodes (16): lockOneRoot(), lockRootRecovering(), relockStaleVaults(), isFileVaultCiphertext(), isRegularFileTarget(), Vault, hasEncryptionPolicy(), isLockedRegularFileErr() (+8 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (76): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+68 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (16): addErrKind, classifyAddErr(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool(), walkEntries() (+8 more)

### Community 9 - "mkdirs"
Cohesion: 0.11
Nodes (42): fileID(), openBunTmp(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir(), TestOpenBunTmp_SymlinkReplacedTargetUntouched(), newWatchPattern(), newFakeInotify() (+34 more)

### Community 10 - "backend.go"
Cohesion: 0.34
Nodes (8): recordTick, isAbsoluteCandidate(), TestCatalogSanity(), parsePasswd(), shell(), TestParsePasswd(), TestParsePasswdMalformed(), hint

### Community 11 - "New"
Cohesion: 0.19
Nodes (19): looksLikeFileVaultRecord(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip(), TestVaultUnlockLockDispatchForRegularFile() (+11 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (19): backendLink, configView, forwardForm, forwardView, protectedPage, recordForm, recordView, trustPage (+11 more)

### Community 13 - "os.File"
Cohesion: 0.04
Nodes (78): multicallAdmissible(), scriptIdentity(), checkFreshDir(), openDirNoFollow(), renewDir(), ReplacementCheck, VettedInode, Multicall (+70 more)

### Community 14 - "plan_test.go"
Cohesion: 0.12
Nodes (25): isWhiteout(), opaque(), pathByte(), pathRecord(), scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent() (+17 more)

### Community 15 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.06
Nodes (19): recordTickCmd(), recordModel, diffModel, editorModel, listenForDaemonEvents(), tickDaemonStats(), listenForNetEvents(), tickNetStats() (+11 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.10
Nodes (8): attrOpCase, eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (57): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+49 more)

### Community 18 - "fyne.io/fyne/v2/widget.Label"
Cohesion: 0.03
Nodes (21): appTheme, binaryRow, clickPanel, counterStrip, homePage, iconTile, iconTileRenderer, navItem (+13 more)

### Community 19 - "fileEditModel"
Cohesion: 0.08
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.06
Nodes (13): ownerPaths(), Resource, innermostResource(), GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase (+5 more)

### Community 23 - "auth.go"
Cohesion: 0.14
Nodes (37): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+29 more)

### Community 24 - "Config"
Cohesion: 0.10
Nodes (31): inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure(), closureRejections(), inGuardedTree() (+23 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.08
Nodes (76): exe_refused(), add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), current_is_inspector() (+68 more)

### Community 26 - "FileTarget"
Cohesion: 0.13
Nodes (10): Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps(), stopPumped() (+2 more)

### Community 27 - "find_test.go"
Cohesion: 0.18
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 28 - "KernelDev"
Cohesion: 0.10
Nodes (27): fdStat, BinaryStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), StatBinary(), StatBinaryFile() (+19 more)

### Community 29 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.12
Nodes (22): placedWalk, treeJudge, confineBelow(), inUserHome(), vouchedFd(), GlobSystemPlaced(), placedDir(), placedEntry() (+14 more)

### Community 30 - "User"
Cohesion: 0.07
Nodes (61): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), setupBunTmpdir(), addedResources(), applyDiffAdditions() (+53 more)

### Community 31 - "TestMonitorUseCaseLifecycle"
Cohesion: 0.18
Nodes (12): NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError(), TestMonitorUseCaseLifecycle(), TestMonitorUseCaseStartError() (+4 more)

### Community 32 - "policyEditor"
Cohesion: 0.14
Nodes (10): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), confirmWarnings() (+2 more)

### Community 33 - "Guard"
Cohesion: 0.06
Nodes (19): deferredBinary, exeRow, rootHandle, supersededBinary, openRoot(), TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), SharedPinDegraded(), TestNoteRootMountFollowsMountID() (+11 more)

### Community 34 - "picker"
Cohesion: 0.09
Nodes (21): keyEntry, listOpts, pickEntry, picker, treeSearch, lstatOf(), outsideHardlinks(), TestSearchTreeFollowsNoSymlink() (+13 more)

### Community 35 - "SelectOrphans"
Cohesion: 0.23
Nodes (16): appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), orphanedPolicies(), orphanedProtectors(), SelectOrphans() (+8 more)

### Community 37 - "Message"
Cohesion: 0.09
Nodes (14): checkSecret(), checkTrustLine(), trustRows(), File, TrustItem, backend, backendOp, backendState (+6 more)

### Community 38 - "controlServer"
Cohesion: 0.11
Nodes (21): readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest(), readGrantRequest() (+13 more)

### Community 39 - "render.go"
Cohesion: 0.07
Nodes (34): lineStyles, viewportState, New(), TestCursorDoesNotBlink(), fg(), Highlighter, lexerFor(), NewHighlighter() (+26 more)

### Community 40 - "pinstate_test.go"
Cohesion: 0.29
Nodes (14): pinOwnerLikelyAlive(), readPinState(), recoverPinState(), pinStateInode(), TestPinOwnerLikelyAlive(), TestPinOwnerLikelyAliveDeadPID(), TestPinStateRoundTrip(), TestReadPinStateMalformed() (+6 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.07
Nodes (15): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), exploitTest, IntegrationSuite (+7 more)

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (8): IntegrationSuite, IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath()

### Community 43 - "monitorUnitTest"
Cohesion: 0.13
Nodes (4): sameInode(), dirTarget(), newMountTable(), monitorUnitTest

### Community 44 - "install_ui.go"
Cohesion: 0.08
Nodes (14): binaryList, installUI, multiAnswer, resizeWatch, sizedDialog, textAnswer, scrollable(), selectable() (+6 more)

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "EventType"
Cohesion: 0.07
Nodes (20): snapEngines, snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting() (+12 more)

### Community 48 - "update/update.go"
Cohesion: 0.04
Nodes (99): installHooks(), staged(), changelogText(), cumulativeChangelog(), newChangelogModel(), newReleaseWindow(), releasesSince(), showChangelog() (+91 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (29): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+21 more)

### Community 50 - "vaultfs_test.go"
Cohesion: 0.12
Nodes (21): DefaultNewFileMode(), IsBinaryContent(), IsOctalMode(), ParseMode(), inodeOf(), mustRead(), mustWrite(), openRootT() (+13 more)

### Community 51 - "editView"
Cohesion: 0.12
Nodes (8): codeEntry, editView, handButton, Entry, fileSize(), joinRel(), newCodeEntry(), relOf()

### Community 52 - "Client"
Cohesion: 0.10
Nodes (11): Client, childEnv(), launch(), Start(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines(), Reader (+3 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.09
Nodes (54): refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), systemPatterns(), systemRule(), watchGroup, addResource(), applyAllowLib() (+46 more)

### Community 54 - "SanitizeText"
Cohesion: 0.05
Nodes (39): newBinaryVetter(), openBinaryVetter(), backingDeviceUnion(), buildGuards(), drainEphemeral(), reloadOnce(), reloadSlotsNeeded(), setConfinedHomes() (+31 more)

### Community 55 - ".apply"
Cohesion: 0.11
Nodes (15): Apply(), DeleteBackup(), Outcome, ProgressFunc, TestApplyKeepsConfirmedLedgerVerdict(), TestApplyPinsUserBinariesAndRetriesStale(), TestApplyRefusesUnpreparedFilesystemBeforeTouchingAnything(), TestApplyRollsBackOnRefusal() (+7 more)

### Community 56 - "go_pkg_strings"
Cohesion: 0.05
Nodes (24): releaseBuild(), Run(), atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape() (+16 more)

### Community 57 - "mutedText"
Cohesion: 0.07
Nodes (44): checkGrid, filesPage, handCheck, recordReview, uninstallCard, newBinaryList(), newBinaryRow(), newCheck() (+36 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (32): markReady(), notifySystemdReady(), engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel() (+24 more)

### Community 59 - "sanitizeTerminalText"
Cohesion: 0.06
Nodes (34): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+26 more)

### Community 60 - "audit.go"
Cohesion: 0.21
Nodes (16): auditAfterEdit(), auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource() (+8 more)

### Community 61 - "Recorder"
Cohesion: 0.15
Nodes (16): confirmExtend(), Target, drain(), Recorder, Target, monitorTargets(), newRecorder(), selfExe() (+8 more)

### Community 62 - "inode_dev"
Cohesion: 0.12
Nodes (23): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free() (+15 more)

### Community 63 - "systemd.go"
Cohesion: 0.05
Nodes (69): anyDaemonRunning(), offerDaemonStart(), runDiffCatalog(), installGUI(), needsBinary(), RunWith(), TestNeedsBinary(), TestRunWithRefusesDeployWithoutBinary() (+61 more)

### Community 64 - "installPage"
Cohesion: 0.15
Nodes (5): actionGate, binSource, installAction, installPage, setWarn()

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm_test.go"
Cohesion: 0.10
Nodes (25): detectWM(), reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket() (+17 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (8): guardBinaryFlag(), infraContainerPaths(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "linkView"
Cohesion: 0.36
Nodes (4): linkTarget(), entryKind, libPlan, linkView

### Community 69 - "time.Duration"
Cohesion: 0.21
Nodes (9): parseErr(), TestLiveSessionConfigAndPut(), TestLiveSessionWatchAndPinger(), holdForward(), pingInterval(), runLiveEdit(), liveSession, RunFileEditorSession() (+1 more)

### Community 70 - "CandidateDir"
Cohesion: 0.08
Nodes (23): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+15 more)

### Community 71 - "RefreshSection"
Cohesion: 0.11
Nodes (25): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+17 more)

### Community 72 - "Candidate"
Cohesion: 0.13
Nodes (25): appendSectionsAndEdit(), collectDiffAdditions(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig(), addManualDirectories(), editConfig(), groupCandidates() (+17 more)

### Community 73 - "fcntl"
Cohesion: 0.08
Nodes (11): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), hold_forever() (+3 more)

### Community 74 - "guardSpec"
Cohesion: 0.21
Nodes (12): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+4 more)

### Community 75 - "planUpdaters"
Cohesion: 0.13
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (18): tempMaskOp, tempRow, TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries() (+10 more)

### Community 78 - "Network"
Cohesion: 0.10
Nodes (13): netGuardPill(), NetworkGuardRepository, newExeResolver(), guardRow(), clonePolicy(), NetEngines, NetExtra, NetPolicy (+5 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "startCatalogRefresh"
Cohesion: 0.08
Nodes (20): swapFile(), newControlManager(), startControlManager(), startControlServer(), makeReloadHandler(), markReadyWithLogin(), startCatalogRefresh(), hasSentinel() (+12 more)

### Community 81 - "op"
Cohesion: 0.15
Nodes (13): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+5 more)

### Community 82 - "LibraryClosure"
Cohesion: 0.09
Nodes (25): defaultLibDirs(), elfInterp(), fileExists(), hostView(), isScript(), LdSoPreloadPaths(), LibraryClosure(), mountVouchesOwnership() (+17 more)

### Community 83 - "rcUser"
Cohesion: 0.13
Nodes (40): rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent(), EnsureBunLauncherEnv(), ensureRCBlock() (+32 more)

### Community 84 - ".FilterExistingWhitelist"
Cohesion: 0.10
Nodes (20): confSafeMatch(), homeMatchConfined(), symlinkStaysInParent(), TestExpandLibDirWritersGlob(), TestFilterExistingWhitelist(), TestFilterExistingWhitelistGlob(), TestFilterExistingWhitelistMixed(), TestFilterExistingWhitelistSkipsOpaqueMulticall() (+12 more)

### Community 85 - "newTestTarget"
Cohesion: 0.09
Nodes (23): NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded(), TestApplyFromWhitelistStopsFirst() (+15 more)

### Community 86 - "runForward"
Cohesion: 0.10
Nodes (23): dialLiveSession(), applyEdit(), editAgain(), fetchConfig(), putConfig(), runEditConfig(), chooseResources(), confirmUserPlaced() (+15 more)

### Community 87 - "updatePage"
Cohesion: 0.21
Nodes (8): installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), boolIndex(), isPreRelease()

### Community 88 - "GenerateConf"
Cohesion: 0.05
Nodes (57): TestDiffMergePreservesExistingAndParses(), collectFilesystemPrereqs(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles() (+49 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.06
Nodes (10): daemonPage, logGrid, logHook, logLine, runHost, runView, logRow(), logStyle() (+2 more)

### Community 92 - "Run"
Cohesion: 0.33
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.16
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.08
Nodes (14): foreignPinnedLink(), PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), judgePasses(), parseProcs() (+6 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "Pulse"
Cohesion: 0.25
Nodes (7): New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse, Sub

### Community 100 - "BinaryEntry"
Cohesion: 0.08
Nodes (40): binaryVerifyState, HashStore, keptEntries, sharedHashEntry, BinaryEntry, checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath() (+32 more)

### Community 101 - "lockAndDeprovision"
Cohesion: 0.14
Nodes (8): deprovisionKind, migrationTarget, classifyDeprovision(), classifyDeprovisionErr(), classifyMigrationTarget(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "recoverFileInPlace"
Cohesion: 0.26
Nodes (15): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), recoverFileInPlace(), rejectSidecarSymlink(), sealFileVaultWithMasterKey(), stageRecovery() (+7 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "Vault"
Cohesion: 0.24
Nodes (3): Vault, Open(), Unopened()

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "openFileVaultWithMasterKey"
Cohesion: 0.19
Nodes (13): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), openFileVaultWithMasterKey(), sealFileVault(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape() (+5 more)

### Community 111 - "cursor.go"
Cohesion: 0.09
Nodes (18): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+10 more)

### Community 112 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.14
Nodes (28): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+20 more)

### Community 113 - "recordFlow"
Cohesion: 0.17
Nodes (11): encryptBlocker(), newRecordFlow(), confirmAndApply(), declined(), offerBackupDeletion(), pickEncryption(), runRecord(), runRecordView() (+3 more)

### Community 115 - "go_pkg_path_filepath"
Cohesion: 0.08
Nodes (44): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+36 more)

### Community 116 - "github.com/spf13/cobra.Command"
Cohesion: 0.25
Nodes (13): credentialFlagState(), ServeConfig, isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestSetupDumpLogRefusesSymlink() (+5 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - ".TestRecord_GuardParity"
Cohesion: 0.35
Nodes (4): IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "runLockdown"
Cohesion: 0.17
Nodes (13): CheckBPFLSM(), configureDaemonLogging(), confirmMasterKeyOverwrite(), resolveConfigPath(), runBPFCheck(), runGenKey(), runLockdown(), runOneShotMode() (+5 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "runMonitor"
Cohesion: 0.21
Nodes (10): CheckEBPF(), MakeDisplayPaths(), newPinGeneration(), prepareDaemonStart(), runHeadless(), runMonitor(), runTUI(), PrintLogo() (+2 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "Monitor"
Cohesion: 0.11
Nodes (3): exeKey, Monitor, watchedInode

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.27
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 131 - "Open"
Cohesion: 0.15
Nodes (21): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+13 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "catalogRefresher"
Cohesion: 0.12
Nodes (14): drain(), newCatalogRefresher(), readEvents(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestDebouncerCoalescesAndBoundsWait() (+6 more)

### Community 143 - "serve.go"
Cohesion: 0.07
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 144 - "NetEventType"
Cohesion: 0.07
Nodes (20): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+12 more)

### Community 145 - "secureResources"
Cohesion: 0.24
Nodes (8): encryptDirectories(), secureResources(), verifyEncryptionState(), resolveFilesystemPrereqs(), askEncryption(), harnessSuite, EnsureSystemSetup(), isKernelAtLeast54()

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (24): candidates(), confirmAll(), describe(), record(), applySchema(), Ledger, Pending, Source (+16 more)

### Community 147 - "runOfflineEdit"
Cohesion: 0.13
Nodes (15): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+7 more)

### Community 148 - "GlobReservations"
Cohesion: 0.17
Nodes (10): GlobChild, statFunc, compileGlobNames(), globText(), GlobReservations, resolveBits(), resolveChildren(), statKey() (+2 more)

### Community 149 - "GuardInodeKey"
Cohesion: 0.05
Nodes (41): LockedError, inodeBatchFn, sweepID, sweepMemo, TestSetInspectorsBounded(), appendInodeRow(), batchInodeRows(), inodeRowsOf() (+33 more)

### Community 150 - "bpfstats/main.go"
Cohesion: 0.38
Nodes (9): chownToSudoUser(), compare(), dumpFullLog(), dumpTrace(), loadOne(), main(), measure(), resizeMaps() (+1 more)

### Community 152 - "Finder"
Cohesion: 0.17
Nodes (9): findButton, span, findAll(), Finder, lowerRunes(), newFindInput(), overlay(), runesEqual() (+1 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.07
Nodes (44): DrainUntilSignal(), buildOneGuard(), eventFilterOptions(), TestEventFilterOptions(), newSelfGuards(), selfProtectSpecs(), startSelfGuard(), startSelfGuards() (+36 more)

### Community 155 - "runNetworkMonitor"
Cohesion: 0.10
Nodes (13): ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), runHeadless(), runNetworkMonitor(), runTUI(), ProtocolString() (+5 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "gui/session.go"
Cohesion: 0.20
Nodes (15): ancestorEnv(), invokingHome(), invokingUID(), ownedBy(), readOwnedFile(), restoreSessionEnv(), statusField(), TestDisplayValue() (+7 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.23
Nodes (22): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+14 more)

### Community 161 - "RawTarget"
Cohesion: 0.31
Nodes (9): BuildEBPFTargets(), RawTarget, isSubDir(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets(), WarnIgnoredFlags() (+1 more)

### Community 162 - "uninstallAll"
Cohesion: 0.08
Nodes (32): deletePostBackups(), listBackups(), restoreBackups(), confirmKeyDeletion(), decryptDirectories(), offerBackupCleanup(), run(), runUninstall() (+24 more)

### Community 163 - "stripTempJournal"
Cohesion: 0.24
Nodes (10): ensurePinStateFilePlaceholder(), ensurePlaceholder(), writeInPlace(), bootstrapSelfProtectPlaceholders(), readTempJournal(), stripTempJournal(), TestTempJournalRoundTripInPlace(), writeTempJournal() (+2 more)

### Community 164 - "NetGuardEvent"
Cohesion: 0.10
Nodes (21): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runHeadless(), runNetworkGuard(), runTUI(), snapNetGuard (+13 more)

### Community 165 - "feed.go"
Cohesion: 0.05
Nodes (43): actionBar, bleedLayout, cell, column, dataTable, eventFeed, feedAction, feedSpec (+35 more)

### Community 166 - "shell"
Cohesion: 0.05
Nodes (56): fakeLink, pageID, procStats, reel, shell, snapNet, statsSample, sysStatus (+48 more)

### Community 167 - "desktop_test.go"
Cohesion: 0.07
Nodes (26): TestEveryInstallModeMapped(), adaptiveGrid, gapVBox, pageDef, stageItem, stageState, tabBar, tabItem (+18 more)

### Community 168 - "watchSet"
Cohesion: 0.19
Nodes (7): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet, Process inspectors

### Community 169 - "AddServeFlags"
Cohesion: 0.36
Nodes (8): AddHeadlessFlag(), AddServeFlags(), TestParseServeFlagsRequiresCredentials(), init(), init(), init(), init(), init()

### Community 170 - "WithSuspectsFrom"
Cohesion: 0.36
Nodes (4): WithSuspectsFrom(), RealNetEngines(), netguardHarness, Option

### Community 172 - "Limitations"
Cohesion: 0.20
Nodes (10): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+2 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "sync.Mutex"
Cohesion: 0.05
Nodes (23): lsmHook, refusedReplacements, tempGrant, trustHook, trustRows, deleteInoKeys(), deleteResKeys(), setSupersedeTrust() (+15 more)

### Community 176 - "safeio.go"
Cohesion: 0.05
Nodes (54): entryOf(), checkAndApply(), TestDaemonUnitWithMetadataOutput(), refreshCatalogConfig(), dupFD(), installConfig(), installFileAs(), installSSHAgent() (+46 more)

### Community 177 - "tempGrantSetup"
Cohesion: 0.25
Nodes (9): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), fakeTempGrant (+1 more)

### Community 178 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 179 - "loadDaemonConfig"
Cohesion: 0.33
Nodes (6): loadDaemonConfig(), noteReloadedSettings(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds()

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 186 - "Binary trust"
Cohesion: 0.25
Nodes (8): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, The dynamic loader, trust

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "Binary"
Cohesion: 0.12
Nodes (24): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), classifyLibs(), finishBinary() (+16 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 197 - "NewDumpHook"
Cohesion: 0.20
Nodes (12): applyVerbosity(), VerboseLevel, BridgeStdLog(), NewDumpHook(), TestBridgeStdLog(), TestDefaultVerboseMirrorsCurrentDisplay(), TestDumpHookFileIsSynced(), TestDumpHookFiltersByVerbosity() (+4 more)

### Community 199 - "Plan"
Cohesion: 0.19
Nodes (8): Section, dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule, Extension

### Community 200 - "NewMonitor"
Cohesion: 0.23
Nodes (6): evalSymlinksOrEmpty(), NewMonitor(), newPathCache(), putBtrfsLayout(), pathCache, pathCacheEntry

### Community 202 - "DaemonEvent"
Cohesion: 0.09
Nodes (39): UIDResolver, NewUIDResolver(), awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemon(), runDaemonUI(), runHeadless() (+31 more)

### Community 203 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Reloading, Runtime sections, Several guarded trees, one vault (+1 more)

### Community 204 - "formatEntry"
Cohesion: 0.50
Nodes (3): formatEntry(), sortedKeys(), sortStrings()

### Community 205 - "ldConfParser"
Cohesion: 0.50
Nodes (3): ldConfParser, ldConfInclude(), ldSoConfDirs()

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "stripImmutableFlags"
Cohesion: 0.50
Nodes (4): stripEntryFlags(), stripImmutableFlags(), TestStripImmutableFlagsDoesNotHangOnFIFO(), TestStripImmutableFlagsSkipsDanglingSymlink()

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "UI"
Cohesion: 0.06
Nodes (32): nextConfig(), InstallMode, InstallRequest, Options, stagedBinary, UninstallHooks, UninstallRequest, UpdateHooks (+24 more)

### Community 213 - "netinject.c"
Cohesion: 0.32
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 216 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 217 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 219 - "TestMain"
Cohesion: 0.50
Nodes (3): TestMain(), helperChild(), TestMain()

### Community 224 - "GuardEvent"
Cohesion: 0.07
Nodes (23): holdHeadless(), HeadlessLine(), bpfGuardEvent, snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), TestLocalizePaths(), commMatchesGuardedBinary() (+15 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 226 - "Prereq"
Cohesion: 0.18
Nodes (7): confirmRunPrereq(), classifySetupError(), TestClassifySetupErrorGeneric(), TestClassifySetupErrorNotSetup(), TestClassifySetupErrorNotSupported(), Prereq, fakeVault

### Community 228 - "mountTable"
Cohesion: 0.19
Nodes (8): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), mountTable

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 241 - "ptrace_seize.c"
Cohesion: 0.50
Nodes (6): main(), peek_buffer(), run_tracer(), run_victim(), touch(), wait_for_file()

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 647 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **35 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `.TestDaemon_Btrfs_CopyIoctlGate`, `.TestRecord_GuardParity`, `IntegrationSuite`, `go_pkg_strings`, `TestMain`?**
  _High betweenness centrality (0.051) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.01925758103868944 - nodes in this community are weakly interconnected._
- **Why does `Editor keys` connect `edit-protected` to `uninstallAll`, `Message`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09125085440874915 - nodes in this community are weakly interconnected._
- **Why does `edit-protected` connect `edit-protected` to `README.md`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **Should `Vault` be split into smaller, more focused modules?**
  _Cohesion score 0.1471264367816092 - nodes in this community are weakly interconnected._